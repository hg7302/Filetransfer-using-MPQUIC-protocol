from  random import uniform
from subprocess import call
import timeit
import os
import json

def dump_to_file(filename,contents):
	f = open(filename,"w")
	f.write(contents)
	f.close()


topo_template = open("topo.template","r").read()
xp_template = open("xp.template","r").read()

schemes = ['quic','https']
losses = [0,2.5]
band_widths = [0,50]
delays = [0,150]
file_sizes = [1024000]
ccs = ['olia','cubic']

number_of_xp = 100
xp_no = 0
exp_run_time_stamp = int(timeit.default_timer())
os.mkdir("/home/mininet/exps/{:d}".format(exp_run_time_stamp))

for _ in range(number_of_xp):
	loss=[uniform(losses[0],losses[1]),uniform(losses[0],losses[1])]
	band_width=[50 + uniform(band_widths[0],band_widths[1]),50 + uniform(band_widths[0],band_widths[1])]
	delay=[uniform(delays[0],delays[1]),uniform(delays[0],delays[1])]
	for file_size in file_sizes:
		for sch in schemes:
			for cc in ccs:
				file_path="/home/mininet/exps/{:d}/exp_{}.json".format(exp_run_time_stamp,xp_no)
				xp_no+=1
				additional=""

				if sch == "quic":
					additional="quicMultipath:{:d}\n".format(1 if cc == "olia" else 0)
				elif cc == "olia":
					call("sudo sysctl -w net.mptcp.mptcp_enabled=1", shell=True)
				else:
					call("sudo sysctl -w net.mptcp.mptcp_enabled=0", shell=True)

				topo_file = topo_template.format(loss=loss,delay=delay,bandwidth=band_width)
				xp_file = xp_template.format(additional=additional,filesize=file_size,type=sch,cc=cc,output=file_path)
				dump_to_file("exp.topo",topo_file)
				dump_to_file("exp.xp",xp_file)
				rc = call("sudo src/mpPerf.py -x exp.xp -t exp.topo >> /dev/null", shell=True)
				f = open(file_path)
				data = json.loads(f.read())
				f.close()
				dura = data['timer']['duration']
				print(sch,cc,file_size,dura,band_width,(int(file_size)/dura/128) * 1.04858)
				print(file_path,rc)
				rc = call("sudo mn -c > /dev/null", shell=True)

				
				