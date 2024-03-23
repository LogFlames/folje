import sacn
import time

sender = sacn.sACNsender()  # provide an IP-Address to bind to if you want to send multicast packets from a specific interface
sender.start()  # start the sending thread
sender.activate_output(7)  # start sending out data in the 1st universe
sender[7].multicast = True  # set multicast to True
sender[7].fps = 30
# sender[1].destination = "192.168.1.20"  # or provide unicast information.
# Keep in mind that if multicast is on, unicast is not used
#start_addr = 361
#d = list(sender[2].dmx_data)

#d[] + 29:start_addr+29+4] = (128, 128, 128, 128)

sender[7].dmx_data = (0, 0, 0, 0)
d = 0
s = time.time()
try:
    while True:
        d = int(50 * (time.time() - s)) % 255
        sender[7].dmx_data = (0,) * 28 + (d, d, d, d)
        sender.flush([7])

except KeyboardInterrupt:
    sender.stop()

