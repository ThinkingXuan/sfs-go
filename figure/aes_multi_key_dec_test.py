# AES-128 AES-192 AES-256 分别200 400 600 800 1000M的文件解密测试

import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")


# 文件大小
users = ["200", "400", "600", "800", "1000"]

# AES128
AES128 = [0.237379333, 0.476022266, 0.7148922, 1.013200266, 1.1921319]
# AES192
AES192 = [0.2595291, 0.520849233, 0.780683433, 1.038239233, 1.306709]
# AES256
AES256 = [0.2925211, 0.58301545, 0.8742694, 1.1697979, 1.46227845]

# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(users))

fig, ax = plt.subplots(figsize=(10, 7))
# AES128
ax.bar(xticks, AES128, width=0.25, label="AES-128", color="y")
# AES192
ax.bar(xticks + 0.25, AES192, width=0.25, label="AES-192", color="green")
# AES256
ax.bar(xticks + 0.5, AES256, width=0.25, label="AES-256", color="red")

# 显示数值
# for a, b in zip(xticks, AES128):
#     ax.text(a, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 
# for a, b in zip(xticks, AES192):
#     ax.text(a+0.25, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 
# for a, b in zip(xticks, AES256):
#     ax.text(a+0.5, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 

ax.set_title("AES Decryption Time", fontsize=15)
ax.set_xlabel("The Size of File(MB)", fontsize=12)
ax.set_ylabel("Elapsed Time(sec)", fontsize=12)
ax.legend()

# 最后调整x轴标签的位置
ax.set_xticks(xticks + 0.25)
ax.set_xticklabels(users)

plt.savefig('tmp.svg',dpi=800, bbox_inches='tight') # 保存成PDF放大后不失真（默认保存在了当前文件夹下）
plt.show()

