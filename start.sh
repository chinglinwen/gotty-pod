
# git="H4sIAAAAAAAA/3RW7W7cuA59lXkB1fcWuED+NneQTdPtomi6wAKLxYCWaIuxTLr6GI/z9Av5K+Nx+mfMcw5JURJH0t/QdcXz8YEYWOMhoyCJTQ/nCVlzih50cwI2XshsSZJwAK8tnbFoUomeMeIbBRyodKjEmZUzohv0Fc10CQFPAf2ZdPYHN0TSQWUG/VbtvICOdMaTBue22jU4lFQbiFBgdKvdoQ/CEHKsPiFHL92wqg1UDSiPhsKhdKIbbYG4GH8VOPTxoKUYEvdI2crxg6oTGcwwRGAD3oQMIoZIXM+8blTqag+TI3HMNWKfQY/YuEF57MTHOVRbFSz4NX5weD1MMlej9uNvizgOVznpl00qHo8/8v6g39LPxyeCB+D6HmSr/IF9+OblBXXcCl8Qu09aS+Ib4RHkMcFPAt7x3+W2ngcHHCxx/bvs/I+fUrTPxy+7iGD33s/Hr/BEsKv/8fgdS4j4NXfGqGjMa11oR7pRFsQAKS3tRnQCrLxwXcJOGxIr6GjD3akh8c7rF8khGYq7FEZ6dgJmS1YqlzJVcquMM/gwDfLhdhBTKdG3zN2H2zJNpd6bzzzuO7R0uKNriNjDsOF+m7hP3z5v514jx1MYQsT21GO5Ee3/3k0frPTmZjJjHjXlmQSSUPxAaP/s7tPwRj0fT/8HekhPtOHuge5B/iJYGjLzj8dTNZ53BC5E8biVcuO5tfHmRJ85JD8ekZv098+n/M/ZkF+f7k9T43rhqJBNUX1EFSJE0lekNWo8RNW4Pm++i6UclYdaHHBdREvejNs3E1raVnhB0yymFiloKVVp3y5kj6WKIi4sxJVWglwoB3TEtWqJaVFaIP74n//eLfiVLok3HuOHB+BLAlbLHfGumO+K94Rd3QvxQly/kMfVY624Ha+OqyIXe7kBZgimfdPWdLdlvtpUgvyC3Vad16lOwC8Eq78jYGMK4rOQxtO4f8grD3QykmvcKXPEzvHQ2a4IDTmnWugmZMHjaBk8zxsZOjL5bzN+xibaMnOLbMlpRWaus+J0WjNV6P2wgNxrZ+qKgTi32Gg78HBGp/JlO2fKtE1SwmTWEC16VYo0Iw42NThaZ+rUsjsZl/l6YpNtVXtJXWHdzyt0eX3x+ZjwEEn4RriChl5tusIp5KmPp8e1V6JXm6dx6cBHRl8sxmG61IsqsESqhgU7qVVHTuJCaOGKatUKU5Q1bNqYGeRpt+CbsBA2tYvZsPQOTb06Tx/FNfFlzTa/I9agtyfVzCwvqs7BWmpN0UGpNKnmLqj8/FiU6C835RvsnLwTeW43gW+PtJv4Tbm9+CZe+Ux9udQ+JeyRXy1y7YiLHVFLjMPWxQNW1KhKfA/5yEvRbhyc1MThrPd5lJM6HCxG4domKkjCaZBkQX4m4H/+BQAA//8BAAD//7hfUhnqCgAA"
# ./gotty-cmd1 -cmd="../gotty-pod" user=wenzhenglin env=online git="$git"
set -e
go build -o gotty-pod1
../gotty/gotty1 --port 7101 -w --permit-arguments ./gotty-pod1 -gitlabtoken=MvPVs7Z56gU2k2ADyR6J