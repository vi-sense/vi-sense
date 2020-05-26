#!/usr/bin/env python3
import sys

if len(sys.argv) != 2:
  print("please specify the file name")
  sys.exit(-1) 

fn = sys.argv[1] 

with open(fn) as f:
  d = set([])
  unique = []
  for i, line in enumerate(f):
    s = line.split(',')
    if len(s) == 3:
      s = line.split(',')
      dt = s[0] + "," + s[1]
      if dt in d:
        print("duplicate on line " + str(i))
      elif s[0] == "" or s[1] == "" or s[2] == "" or len(s) != 3:
        print("line not complete")
      else:
        d.add(dt)
        unique.append(line)


  s = fn.split('.')
  with open(s[0] + "-unique.csv", 'w') as f:
      for item in unique:
          f.write("%s" % item)
