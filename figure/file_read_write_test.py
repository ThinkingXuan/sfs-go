
# 个人文件读写性能测试
from matplotlib import ticker
import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")


# 文件大小
users = [5, 10, 20, 50, 80, 100, 200, 300, 400, 500, 600, 700, 800, 900,1000]

# 上传时间        5     10      20    50     80     100    200    300    400    500    600    700     800   900    1000
upload_time = [0.361, 0.425, 0.539, 0.874, 1.303, 1.468, 2.540, 3.536, 4.447, 5.568, 6.205, 6.845 , 7.571, 8.534 ,10.455]

# 下载时间          5     10      20     50     80    100    200    300    400     500    600     700      800    900    1000
download_time = [0.125, 0.282, 0.301, 0.833, 1.743, 2.308, 3.60, 4.999,  5.804, 7.057, 7.559, 8.959,  11.054, 13.254, 16.851]


# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(users))
print(xticks) 

fig, ax = plt.subplots(figsize=(10, 7))
# 所有门店第一种产品的销量，注意控制柱子的宽度，这里选择0.25
ax.plot(users, upload_time, color='green', marker='o', linewidth=1.5, label="File Upload")#DE6B58

ax.plot(users, download_time, color='#DE6B58', marker='o', linewidth=2, label="File Download")

# 文件上传显示数值
for a, b in zip(users, upload_time):
    if a < 100:
        continue
    ax.text(a, b, '%.3f'%b, ha='center', va= 'bottom',fontsize=12) 
 # 文件下载显示数值
 
for a, b in zip(users, download_time):
    if a < 100:
        continue
    ax.text(a, b, '%.3f'%b, ha='center', va= 'bottom',fontsize=12) 


ax.set_title("BFShare's File Upload and Download Time", fontsize=15)
ax.set_xlabel("The Size of File(MB)", fontsize=12)
ax.set_ylabel("Elapsed Time(sec)", fontsize=12)
ax.legend()

# 最后调整x轴标签的位置
new_ticks = [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000]
ax.set_xticks(new_ticks)

ax.set_xticklabels(new_ticks)

plt.savefig('tmp.svg',dpi=800, bbox_inches='tight') 
plt.show()