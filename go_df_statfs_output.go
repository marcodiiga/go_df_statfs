package main

import (
	"fmt"
	"log"
	"os/exec"

	"golang.org/x/sys/unix"
)

func main() {

	// Print df / and df -i / for comparison
	fmt.Println("--- df / ---")
	out, err := exec.Command("df", "/").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
	fmt.Println("--- df -i / ---")
	out, err = exec.Command("df", "-i", "/").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
	fmt.Println("--- statfs ---")

	var statFsOutput unix.Statfs_t
	err = unix.Statfs("/", &statFsOutput)
	if err != nil {
		log.Fatal(err)
	}

	// Not available, see note below
	//
	// var statVfsOutput unix.Statvfs_t
	// err = unix.Statvfs("/", &statVfsOutput)
	// if err != nil {
	// 	fmt.Printf("error: %v", err)
	// }

	totalBlockCount := statFsOutput.Blocks     // Total data blocks in filesystem
	freeSystemBlockCount := statFsOutput.Bfree // Free blocks in filesystem for root - rocket-s3 doesn't run as root so no need to concern with this
	freeUserBlockCount := statFsOutput.Bavail  // Free blocks available to unprivileged user

	// The total count of  blocks available to unprivileged users is not the totalBlockCount but rather, assuming
	// extra root space is always filled _after_ user space has been filled, the total number of blocks available
	// _minus_ that extra root-only space.
	totalBlockCountAvailableToUnprivilegedUsers := totalBlockCount - (freeSystemBlockCount - freeUserBlockCount)
	if freeSystemBlockCount > totalBlockCount { // Might happen on some filesystems, we're just interested in user space though
		totalBlockCountAvailableToUnprivilegedUsers = totalBlockCount
	}
	usedUserBlockCount := totalBlockCountAvailableToUnprivilegedUsers - freeUserBlockCount

	// Assumes 4K blocks
	fmt.Printf("Frsize: %d, Bsize: %d\n", statFsOutput.Frsize, statFsOutput.Bsize)
	fmt.Printf("Total block count: %d, total count of 1-K blocks (df / output): %d\n", totalBlockCount, totalBlockCount*4)
	fmt.Printf("Used user block count: %d, used 1-K user blocks (df / output): %d\n", usedUserBlockCount, usedUserBlockCount*4)
	fmt.Printf("Free user block count: %d, free 1-K user blocks (df / output): %d\n", freeUserBlockCount, freeUserBlockCount*4)
	fmt.Printf("percentage of used space: %d %%\n", 100-(100*freeUserBlockCount/totalBlockCountAvailableToUnprivilegedUsers))

	totalInodeCount := statFsOutput.Files
	freeSystemInodeCount := statFsOutput.Ffree
	// NOTE: statfs is a non-posix compliant version of statvfs (even though sometimes it can be even more portable..)
	// it doesn't return a structure which distinguishes between 'free and available inodes', i.e. inodes free for root
	// and inodes free for non-root users, it just uses the 'free inodes' field (for root). To be even more accurate in our
	// usecase we should use statvfs, but go extended unix branch doesn't provide it. Unless we want to do the hardlinking
	// ourselves (and in cgo case it's a solid _NOPE_ for just this one thing), let's just ignore it for the moment being.

	fmt.Printf("Total inode count (df -i / output): %d\n", totalInodeCount)
	fmt.Printf("Free system inode count (df / output): %d\n", freeSystemInodeCount)
	fmt.Printf("percentage of used inodes: %d %%\n", 100-(100*freeSystemInodeCount/totalInodeCount))
}
