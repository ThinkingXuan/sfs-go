# 动态混合加密测试，分别200 400 600 800 1000M的文件加密和解密测试

import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")

# 文件大小
filesize = ["200", "400", "600", "800", "1000"]

dec = [0.3595291, 0.620849233, 0.880683433, 1.238239233, 1.356709]


# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(filesize))

fig, ax = plt.subplots(figsize=(10, 7))
# # AES128
# ax.bar(xticks, AES128, width=0.25, label="AES-128", color="y")
ax.bar(xticks + 0.25, dec, width=0.25, label="Hybrid Decryption", color="red")

# 显示数值
# for a, b in zip(xticks, AES128):
#     ax.text(a, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 
# for a, b in zip(xticks, AES192):
#     ax.text(a+0.25, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 
# for a, b in zip(xticks, AES256):
#     ax.text(a+0.5, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 

ax.set_title("Dynamic Hybrid Decryption Time", fontsize=15)
ax.set_xlabel("The Size of File(MB)", fontsize=12)
ax.set_ylabel("Elapsed Time(sec)", fontsize=12)
ax.legend()

# 最后调整x轴标签的位置
ax.set_xticks(xticks + 0.25)
ax.set_xticklabels(filesize)

plt.savefig('tmp.svg',dpi=800, bbox_inches='tight') # 保存成PDF放大后不失真（默认保存在了当前文件夹下）
plt.show()

