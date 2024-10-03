### Начало анализа
Сперва я просто попробовал создать таблицы, с которыми будет просто работать.

Вот так выглядела таблица, где хранились заказы пользователей.
```
create table if not exists orders (
    user_id bigint not null,
    order_id bigint not null,
    expiration_date text not null,
    package_type text not null,
    weight bigint not null,
    cost bigint not null,
    use_tape boolean not null,
    primary key(user_id, order_id)
);
```

Таблица с историей:
```
create table if not exists orders_history (
    user_id bigint not null,
    order_id bigint not null,
    expiration_date text not null,
    package_type text not null,
    weight bigint not null,
    cost bigint not null,
    use_tape boolean not null,
    status text not null,
    updated_at text not null,
    primary key(order_id)
);
```

Таблица с возвратами:
```
create table if not exists refunds (
    order_id bigint not null,
    primary key(order_id)
);
```

Сразу же можно заметить дублирование данных в 2 таблицах, что не есть хорошо, но поначалу меня это не беспокоило.

### Анализ запроса к orders
Перед анализом запросов я добавил в бд 1000 заказов, 500 из них "пользователи" забрали, 250 вернули и 125 заказов вернул курьеру. Этого я достиг с помощью запуска интеграционного теста. Единственное я не удалял контейнер после теста, чтобы можно было проанализировать запросы.

Для анализа на данном этапе я выбрал следующий запрос, который ищет заказ по id пользователя и заказа:

```
explain (analyze, verbose, buffers)
select user_id,
    order_id,
    expiration_date,
    package_type,
    cost,
    weight,
    use_tape
from orders
where user_id = 4570190542752640735
    and order_id = 7081787540312482718
```

Получили следующее:

```
Index Scan using orders_pkey on public.orders  
(cost=0.27..8.29 rows=1 width=50) (actual time=0.044..0.045 rows=1 loops=1)
Output: user_id, order_id, expiration_date, package_type, cost, weight, use_tape
Index Cond: ((orders.user_id = '4570190542752640735'::bigint) AND
            (orders.order_id = '7081787540312482718'::bigint))
Buffers: shared hit=3
Planning:
    Buffers: shared hit=86
    Planning Time: 1.136 ms
    Execution Time: 0.116 ms
```

Использовался Index Scan по primary keys таблицы order. 

Добавим индекс по order_id и user_id и сравним результаты:
```
Index Scan using orders_user_id_order_id_idx on public.orders 
(cost=0.27..8.29 rows=1 width=50) (actual time=0.021..0.021 rows=1 loops=1)
Output: user_id, order_id, expiration_date, package_type, cost, weight, use_tape
Index Cond: ((orders.user_id = '7784818529561224137'::bigint) AND 
            (orders.order_id = '3277531223011217332'::bigint))
Buffers: shared hit=3
Planning:
    Buffers: shared hit=102
    Planning Time: 0.660 ms
    Execution Time: 0.055 ms
```

В этот раз использовались созданные индексы. Время уменьшилось в 2 раза.

### Анализ запроса к orders_history

Теперь будем искать заказ в orders_history по его id:

```
explain (analyze, verbose, buffers)
select user_id,
    expiration_date,
    package_type,
    weight,
    cost,
    use_tape,
    status,
    updated_at
from orders_history
where order_id = 6256450497015182523
```

Получаем:

```
Index Scan using orders_history_pkey on public.orders_history  
(cost=0.28..8.29 rows=1 width=68) (actual time=0.025..0.026 rows=1 loops=1)
Output: user_id, expiration_date, package_type, weight, cost, use_tape, status, updated_at
Index Cond: (orders_history.order_id = '6256450497015182523'::bigint)
Buffers: shared hit=3
Planning:
    Buffers: shared hit=86
    Planning Time: 0.556 ms
    Execution Time: 0.074 ms
```

Добавим индекс и сравним результаты:

```
Index Scan using orders_history_order_id_idx on public.orders_history  
(cost=0.28..8.29 rows=1 width=68) (actual time=0.024..0.024 rows=1 loops=1)
Output: user_id, expiration_date, package_type, weight, cost, use_tape, status, updated_at
Index Cond: (orders_history.order_id = '7741481755526607455'::bigint)
Buffers: shared hit=3
Planning:
    Buffers: shared hit=99
    Planning Time: 0.621 ms
    Execution Time: 0.074 ms
```

В данном случае особой разницы нет, но может быть при бОльших размерах бд прирост производительности будет выше.

### Анализ запроса к refunds + orders_history
В refunds я храню только номера заказов. Для того, чтобы получить список возращенных заказов с подробной информацией нужно делать запрос сразу к 2 таблицам.

Вот сам запрос:

```
select
	oh.user_id,
	oh.order_id,
	oh.expiration_date,
	oh.package_type,
	oh.weight,
	oh.cost,
	oh.use_tape
	from orders_history oh
	join (
		select order_id
		from refunds
		order by order_id
			limit 50 offset 0
	) r on oh.order_id = r.order_id
	order by oh.order_id
```

Получился очень страшный и длинный вывод:

```
Sort  (cost=30.46..30.59 rows=50 width=97) 
(actual time=0.276..0.279 rows=50 loops=1)"
  Output: oh.user_id, oh.order_id, oh.expiration_date, oh.package_type, oh.weight, oh.cost, oh.use_tape
  Sort Key: oh.order_id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=34
  ->  Hash Join  (cost=3.09..29.05 rows=50 width=97) 
      (actual time=0.098..0.231 rows=50 loops=1)
        Output: oh.user_id, oh.order_id, oh.expiration_date, oh.package_type, oh.weight, oh.cost, oh.use_tape
        Hash Cond: (oh.order_id = refunds.order_id)
        Buffers: shared hit=31
        ->  Seq Scan on public.orders_history oh  (cost=0.00..22.88 rows=688 width=97)
            (actual time=0.004..0.089 rows=1000 loops=1)
              Output: oh.user_id, oh.order_id, oh.expiration_date, oh.package_type, oh.weight, oh.cost, oh.use_tape, oh.status, oh.updated_at
              Buffers: shared hit=16
        ->  Hash  (cost=2.47..2.47 rows=50 width=8) (actual time=0.064..0.064 rows=50 loops=1)
              Output: refunds.order_id
              Buckets: 1024  Batches: 1  Memory Usage: 10kB
              Buffers: shared hit=15
              ->  Limit  (cost=0.15..1.97 rows=50 width=8) (actual time=0.026..0.052 rows=50 loops=1)
                    Output: refunds.order_id
                    Buffers: shared hit=15
                    ->  Index Only Scan using refunds_pkey on public.refunds  (cost=0.15..82.06 rows=2260 width=8) (actual time=0.025..0.048 rows=50 loops=1)
                          Output: refunds.order_id
                          Heap Fetches: 83
                          Buffers: shared hit=15
Planning:
  Buffers: shared hit=102
Planning Time: 0.821 ms
Execution Time: 0.345 ms
```

Далее я попробовал создать игдекс на order_id в refunds, но практически ничего не изменилось при использование такого запроса.
Поэтому я увеличил количество записей в бд. В orders_history 50k записей, в refunds 6k.

Далее я повторил тот же запрос:

```
Nested Loop  (cost=0.57..390.48 rows=50 width=53) (actual time=0.066..0.358 rows=50 loops=1)
  Output: oh.user_id, oh.order_id, oh.expiration_date, oh.package_type, oh.weight, oh.cost, oh.use_tape
  Inner Unique: true
  Buffers: shared hit=200
  ->  Limit  (cost=0.28..2.61 rows=50 width=8) (actual time=0.038..0.092 rows=50 loops=1)
        Output: refunds.order_id
        Buffers: shared hit=50
        ->  Index Only Scan using refunds_order_id_idx on public.refunds 
            (cost=0.28..283.81 rows=6102 width=8) (actual time=0.037..0.088 rows=50 loops=1)
              Output: refunds.order_id
              Heap Fetches: 50
              Buffers: shared hit=50
  ->  Index Scan using orders_history_order_id_idx on public.orders_history oh 
      (cost=0.29..7.75 rows=1 width=53) (actual time=0.005..0.005 rows=1 loops=50)
        Output: oh.user_id, oh.order_id, oh.expiration_date, oh.package_type, oh.weight, oh.cost, oh.use_tape, oh.status, oh.updated_at
        Index Cond: (oh.order_id = refunds.order_id)
        Buffers: shared hit=150
Planning:
  Buffers: shared hit=131
Planning Time: 0.888 ms
Execution Time: 0.413 ms
``` 

Теперь используются индексы + Nested Loop, т.к. limit указан 50 (не так много). Если же увеличить limit и offset, то будут использоваться индексы, Hash Join и Sort.

### Оптимизация структуры таблиц

Как уже упоминалось в самом начале, данные дублируются в двух таблицах, поэтому я перепишу таблицы и запросы к ним так, чтобы основная информация по заказу хранилась только в истории заказа.

Так теперь выглядит таблица orders:

```
create table if not exists orders (
    user_id bigint not null,
    order_id bigint not null,
    primary key(user_id, order_id)
);
```

Также для user_id в таблице orders_history создал индекс.