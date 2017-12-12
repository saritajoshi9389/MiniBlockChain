# Automation script for creating the server config
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
import json
import sys
my_dict = {}
count = int(sys.argv[1])
port = 5000
my_dict["count"] = count
# Open a file for writing
for i in range(1, count+1):
    my_dict[i] = {
        "server_ip": "localhost",
        "server_port": str(5000 + i),
        "data_directory_temp": "/tmp/mini_blockchain/"+str(i)+"/"

    }
out_file = open("config."+str(count)+".json","w")
json.dump(my_dict,out_file, indent=4)
# Close the file
out_file.close()