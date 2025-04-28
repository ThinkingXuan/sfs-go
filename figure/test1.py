## 论文绘制图
# encoding=utf-8
import matplotlib.pyplot as plt
from scipy.interpolate import make_interp_spline

names = ['0', '5','10','20','50','100','200','500','800','1000']
x = range(len(names))
print(x)
y = [0.0, 0.294, 0.364, 0.541, 1.02, 1.395, 2.926, 5.805, 8.961, 12.085]
# y1=[0.86,0.85,0.853,0.849,0.83,1]
#plt.plot(x, y, 'ro-')
#plt.plot(x, y1, 'bo-')
#pl.xlim(-1, 11)  # 限定横轴的范围
#pl.ylim(-1, 110)  # 限定纵轴的范围
plt.plot(x, y, marker='o', mec='r', mfc='w',label=u'SFSharing file upload')
# plt.plot(x, y1, marker='*', ms=10,label=u'y=x^3曲线图')
plt.legend()  # 让图例生效
plt.xticks(x, names)
plt.margins(0)
plt.subplots_adjust(bottom=0.15)

plt.xlabel(u"File Size(MB)") #X轴标签
plt.ylabel("Time(sec)") #Y轴标签

# plt.title("A simple plot") #标题

plt.show()
