package strategy

import (
	"fmt"
)

const (
	CostTape    = 1
	CostPackage = 5
	CostBox     = 20

	// Вес в граммах
	PackageMaxWeight = 10 * 1000
	BoxMaxWeight     = 30 * 1000
)

var ContainerTypeMap = map[string]ContainerStrategy{
	"":        &DefaultContainerStrategy{},
	"package": &PackageStrategy{},
	"box":     &BoxStrategy{},
	"tape":    &TapeStrategy{},
}

type ContainerStrategy interface {
	Type() string
	UseTape() error
	IsTaped() bool
	CalculateCost(weight, cost uint64) (uint64, error)
}

type DefaultContainerStrategy struct{}

func (s *DefaultContainerStrategy) Type() string {
	return "default"
}

func (s *DefaultContainerStrategy) UseTape() error {
	return nil
}

func (s *DefaultContainerStrategy) IsTaped() bool {
	return false
}

func (s *DefaultContainerStrategy) CalculateCost(weight, cost uint64) (uint64, error) {
	return cost, nil
}

type PackageStrategy struct {
	useTape bool
}

func (s *PackageStrategy) Type() string {
	if s.useTape {
		return "taped package"
	}

	return "package"
}

func (s *PackageStrategy) UseTape() error {
	s.useTape = true
	return nil
}

func (s *PackageStrategy) IsTaped() bool {
	return s.useTape
}

func (s *PackageStrategy) CalculateCost(weight, cost uint64) (uint64, error) {
	if weight > PackageMaxWeight {
		return 0, fmt.Errorf("max weight for package is %dgr", PackageMaxWeight)
	}

	res_cost := cost + CostPackage
	if s.useTape {
		res_cost += CostTape
	}

	return res_cost, nil
}

type BoxStrategy struct {
	useTape bool
}

func (s *BoxStrategy) Type() string {
	if s.useTape {
		return "taped box"
	}

	return "box"
}

func (s *BoxStrategy) UseTape() error {
	s.useTape = true
	return nil
}

func (s *BoxStrategy) IsTaped() bool {
	return s.useTape
}

func (s *BoxStrategy) CalculateCost(weight, cost uint64) (uint64, error) {
	if weight > BoxMaxWeight {
		return 0, fmt.Errorf("max weight for box is %dgr", BoxMaxWeight)
	}

	res_cost := cost + CostBox
	if s.useTape {
		res_cost += CostTape
	}

	return res_cost, nil
}

type TapeStrategy struct{}

func (s *TapeStrategy) Type() string {
	return "tape"
}

func (s *TapeStrategy) UseTape() error {
	return fmt.Errorf("can't use tape twice")
}

func (s *TapeStrategy) IsTaped() bool {
	return true
}

func (s *TapeStrategy) CalculateCost(weight, cost uint64) (uint64, error) {
	return cost + CostTape, nil
}
