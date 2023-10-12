package main

import (
	"fmt"
	"linux/blockchainConnection"
	"linux/memory"
	"linux/storage"
	"linux/uptime"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shirou/gopsutil/cpu"
)

type totalCalculationStruct struct {
	cpuCores uint64
	totalRAM float64
	upTime   string
	hdTotal  float64
	cpuUsage float64
}

func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	fmt.Printf("CPU cores: %d\n", memory.GetCPUCoreCount())
	fmt.Printf("Total RAM: %.2f GB\n", memory.GetTotalRAM())
	uptime, err := uptime.GetUptime()
	if err != nil {
		fmt.Println("Error getting system uptime:", err)
		return
	}
	fmt.Println("System Uptime:", uptime.String())

	disks := storage.GetHardDisks()
	if err != nil {
		fmt.Println("Error getting hard disks:", err)
		return
	}
	var totalHDiskSize float64
	fmt.Printf("%-10s                %-20s%-20s\n", "DeviceID", "Size", "Free")
	for _, disk := range disks {
		totalHDiskSize = disk.SizeBytes + disk.FreeBytes
		fmt.Printf("%-10s          %-20s\n", disk.Name, totalHDiskSize)
	}

	percent, err := cpu.Percent(time.Second, true)
	if err != nil {
		fmt.Println("Error getting CPU usage:", err)
		return
	}

	fmt.Printf("CPU usage: %.2f%%\n", percent[0])

	totalcal := totalCalculationStruct{
		cpuCores: calculateCORE(memory.GetCPUCoreCount()),
		totalRAM: calculateRAM(memory.GetTotalRAM()),
		upTime:   uptime.String(),
		hdTotal:  calculateSSDS(totalHDiskSize),
		cpuUsage: percent[0],
	}

	result := totalCalculation(totalcal)
	fmt.Println("this is comming final result", result)

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// // systemSleep()

	// <-interrupt // Wait for a signal to exit

	// fmt.Println("Exiting...")
}

type stackESG struct {
	nodeStake uint64
	esgScore  uint64
}

// each validator node stacking and esg score get

func validatorNodeStaking() stackESG {
	validatorAddress := common.HexToAddress("0xD4De91e65770fCd8b58044a539ff338A3938fD70")
	result := blockchainConnection.GetNodeStakeWithEsg(validatorAddress)
	// Manually extract the values from the result
	stakedAmount := new(big.Int).SetBytes(result[0:32])
	esg := new(big.Int).SetBytes(result[32:64])
	// lockTimestamp := new(big.Int).SetBytes(result[64:96])
	return stackESG{nodeStake: uint64(stakedAmount.Int64()), esgScore: uint64(esg.Int64())}
}

// total stack of all node on validator nodes

func totalStake() *big.Int {
	result := blockchainConnection.GettotalStack()
	stakedValue := new(big.Int)
	stakedValue.SetBytes(result)
	return stakedValue
}

// get all validator list

func validatorListAsArray() []common.Address {
	result := blockchainConnection.GetValidatorList()
	var validatorList []common.Address
	for i := 0; i < len(result); i += 32 {
		addressBytes := result[i : i+32]
		validatorList = append(validatorList, common.BytesToAddress(addressBytes))
	}

	return validatorList
}

func totalCalculation(TC totalCalculationStruct) float64 {
	// total stack come from
	totalStack := totalStake()
	// node one stack
	nodeStake := validatorNodeStaking()
	eachNodeStake := nodeStake.nodeStake
	SCi := calculateStakeCoefficient(uint64(totalStack.Int64()), eachNodeStake)
	// EsG score from smart contract
	esgScore := float64(nodeStake.esgScore)
	//total uptime calculation from node
	epoch := 15 //uptime
	totalUptime := float64((15 - 10) / epoch)

	ssdv := calculateSSDS(TC.hdTotal)
	ramv := calculateRAM(TC.totalRAM)
	cpuutilv := calculateCORE(TC.cpuCores)
	// USi computed as (100 – total blocks mined since last uptime) / 100 where USi has a lower bound of 0 ie at any point
	// TODO when blockchain will be run that time i will calulate proper mined blocks of each node in last uptime
	totalBlockMined := 100
	nodeUtilization := float64(100-totalBlockMined) / 100
	nodecapacity := nodeCapacityScore(ssdv, ramv, cpuutilv)
	fmt.Println("values - SCi -totalUptime -nodecapacity -esgScore -nodeUtilization", SCi, totalUptime, nodecapacity, esgScore, nodeUtilization)
	return (SCi * 0.3) + (totalUptime * 0.2) + (nodecapacity * 0.15) + (esgScore / 100 * 0.25) + (nodeUtilization * 0.1)
}

func calculateSSDS(solidStateDriveCapacity float64) float64 {
	// Apply constraints for minimum and maximum scores
	var ssds float64
	switch {
	case solidStateDriveCapacity <= 2048:
		ssds = 0.0
	case solidStateDriveCapacity >= 16384:
		ssds = 1.0
	default:
		ssds = (solidStateDriveCapacity / (1024 * 14.0)) - (1.0 / 7.0)
	}

	return ssds
}

func calculateRAM(ramInstalled float64) float64 {
	// Apply constraints for minimum and maximum score
	var ram float64
	switch {
	case ramInstalled <= 4:
		ram = 00
	case ramInstalled >= 32:
		ram = 1.0
	default:
		ram = (ramInstalled / 28) - (1.0 / 7.0)
	}

	return ram
}

func calculateCORE(coresInstalled uint64) uint64 {
	// Apply constraints for minimum and maximum scores
	var core uint64
	switch {
	case coresInstalled <= 4:
		core = 0.0
	case coresInstalled >= 16:
		core = 1.0
	default:
		core = (uint64(coresInstalled) / 12) - (1 / 3)
	}
	return core
}

func calculateStakeCoefficient(S uint64, stakedAmounts uint64) float64 {
	return float64(S) / float64(stakedAmounts)
}

func nodeCapacityScore(ssdv float64, ramv float64, cpuutilv uint64) float64 {
	ssdScore := calculateSSDS(ssdv) * 0.2
	ramScore := calculateRAM(ramv) * 0.4
	cpuutilScore := float64(calculateCORE(cpuutilv)) * 0.4 // put cpuutilization here
	nodeCapacityScore := ssdScore + ramScore + cpuutilScore
	return nodeCapacityScore
}
