package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func TestAddOrder(t *testing.T) {
	expected_output := ""
	storagePATH := "storage_e2e_add_order.json"
	defer os.Remove(storagePATH)

	expDate := utils.CurrentDateString()
	args := []string{"accept", "order", "-u", "432", "-o", "31", "-c", "1402", "-w", "402", "-p", "package", "-s", "-t", expDate}

	cmdName := "../bin/manager"
	cmd := exec.Command(cmdName, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "STORAGE_PATH="+storagePATH)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, expected_output, string(output))
}

func TestAddOrderWrongPackageType(t *testing.T) {
	expected_output := "wrong isn't container type\n"
	storagePATH := "storage_e2e_add_order.json"
	defer os.Remove(storagePATH)

	expDate := utils.CurrentDateString()
	args := []string{"accept", "order", "-u", "432", "-o", "31", "-c", "1402", "-w", "402", "-p", "wrong", "-s", "-t", expDate}

	cmdName := "../bin/manager"
	cmd := exec.Command(cmdName, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "STORAGE_PATH="+storagePATH)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, expected_output, string(output))
}

func TestGiveOrder(t *testing.T) {
	expected_output := ""
	storagePATH := "storage_e2e_give_order.json"
	defer os.Remove(storagePATH)

	expDate := utils.CurrentDateString()
	cmdName := "../bin/manager"
	cmd := exec.Command(cmdName, "accept", "order", "-u", "1312", "-o", "10", "-c", "3420", "-w", "3234", "-p", "package", "-s", "-t", expDate)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "STORAGE_PATH="+storagePATH)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, expected_output, string(output))

	cmd = exec.Command(cmdName, "give", "-o=10")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "STORAGE_PATH="+storagePATH)

	output, err = cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, expected_output, string(output))
}
