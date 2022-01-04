```
$ go run ./go_df_statfs_output.go 
--- df / ---
Filesystem                1K-blocks      Used Available Use% Mounted on
/dev/mapper/vgubuntu-root 515010816 202741436 286038600  42% /

--- df -i / ---
Filesystem                  Inodes   IUsed    IFree IUse% Mounted on
/dev/mapper/vgubuntu-root 32768000 3111871 29656129   10% /

--- statfs ---
Frsize: 4096, Bsize: 4096
Total block count: 128752704, total count of 1-K blocks (df / output): 515010816
Used user block count: 50685359, used 1-K user blocks (df / output): 202741436
Free user block count: 71509650, free 1-K user blocks (df / output): 286038600
percentage of used space: 42 %
Total inode count (df -i / output): 32768000
Free system inode count (df / output): 29656129
percentage of used inodes: 10 %
```
