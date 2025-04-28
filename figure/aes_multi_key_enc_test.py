# AES-128 AES-192 AES-256 分别200 400 600 800 1000M的文件加密测试

import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")


# 文件大小
users = ["200", "400", "600", "800", "1000"]

# AES128
AES128 = [0.29560, 0.589693, 0.899511, 1.179075166, 1.475104766]
# AES192
AES192 = [0.330859466, 0.6301771, 0.917658666, 1.2780507, 1.591336366]
# AES256
AES256 = [0.35910085, 0.7095105, 1.0401628, 1.4204446, 1.7813629]

# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(users))

fig, ax = plt.subplots(figsize=(10, 7))
# AES128
ax.bar(xticks, AES128, width=0.25, label="AES-128", color="y")
# AES192
ax.bar(xticks + 0.25, AES192, width=0.25, label="AES-192", color="green")
# AES256
ax.bar(xticks + 0.5, AES256, width=0.25, label="AES-256", color="red")

# # 显示数值
# for a, b in zip(xticks, AES128):
#     ax.text(a, b+1, '%.2f'%b, ha='center', va= 'bottom',fontsize=10) 
# for a, b in zip(xticks, AES192):
#     ax.text(a+0.25, b+1, '%.2f'%b, ha='center', va= 'bottom',fontsize=10) 
# for a, b in zip(xticks, AES256):
#     ax.text(a+0.5, b+1, '%.2f'%b, ha='center', va= 'bottom',fontsize=10) 

ax.set_title("AES Encryption Time", fontsize=15)
ax.set_xlabel("The Size of File(MB)", fontsize=12)
ax.set_ylabel("Elapsed Time(sec)", fontsize=12)
ax.legend()

# 最后调整x轴标签的位置
ax.set_xticks(xticks + 0.25)
ax.set_xticklabels(users)

plt.savefig('tmp.svg',dpi=800, bbox_inches='tight') # 保存成PDF放大后不失真（默认保存在了当前文件夹下）
plt.show()

