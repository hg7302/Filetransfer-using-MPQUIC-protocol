import json

f = open('files','r')
files = f.read()
f.close()
items = []

def get_loss(a):
	return float(a.split()[1][:-1]);

def get_bw_and_delay(a):
	a = a.split(',')
	return float(a[2]),float(a[0])

for exp_file in files.split():
	print(exp_file)
	line = open(exp_file,'r').read()
	if len(line) == 0:
		continue
	data = json.loads(line)
	flat_dict = {}
	for (_,value) in data.items():
		for (k,v) in value.items():
			flat_dict[k] = v
	query_dict = {}
	query_dict['loss_0'] = get_loss(flat_dict['netemAt_0'])
	query_dict['loss_1'] = get_loss(flat_dict['netemAt_1'])
	
	query_dict['bw_0'],query_dict['delay_0'] = get_bw_and_delay(flat_dict['path_0'])
	query_dict['bw_1'],query_dict['delay_1'] = get_bw_and_delay(flat_dict['path_1'])

	query_dict['file_size'] = int(flat_dict['file_size'])
	query_dict['duration'] = flat_dict['duration']
	query_dict['outputPath'] = flat_dict['outputPath']
	query_dict['xpType'] = flat_dict['xpType'] + ('_MP' if flat_dict['congctrl'] == 'olia' else '')
	query_dict['goodput'] = ((query_dict['file_size']* 1.04858 /128)/(query_dict['duration']))
	items.append(query_dict)


create_table = """
CREATE TABLE IF NOT EXISTS
exp_analytics(
	xpType VARCHAR(255),
	file_size INTEGER,
	duration FLOAT,
	loss_0 FLOAT,
	loss_1 FLOAT,
	delay_0 FLOAT,
	delay_1 FLOAT,
	bw_0 FLOAT,
	bw_1 FLOAT,
	goodput FLOAT,
	outputPath VARCHAR(255)
)
"""

st = ""
for item in items:
	# print(item['bw_1'],item['bw_0'],item['file_size'],item['duration'],item['goodput'],item['xpType'])
	st = st + "('{d[xpType]}',{d[file_size]:d},{d[duration]:f},{d[loss_0]:f},{d[loss_1]:f},{d[delay_0]:f},{d[delay_1]:f},{d[bw_0]:f},{d[bw_1]:f},{d[goodput]:f},'{d[outputPath]}')\n".format(d=item)
	st += ','

print(st[:-1])
