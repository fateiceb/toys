srv_port1 = 8080
srv_port2 = 8888
lsof -i:${srv_port1} | awk '{if (NR == 2) {print $2}}' | xargs kill -9
lsof -i:${srv_port2} | awk '{if (NR == 2) {print $2}}' | xargs kill -9
