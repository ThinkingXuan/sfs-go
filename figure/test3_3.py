# 实验 分享10个文件的落块数量 折线图

import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")


# 用户数量
users = ["2", "4", "6", "8", "10", "12"]

# 总落块量
all_block = [49, 90, 128, 168, 202, 231]

# 总落块量
all_block2 = [194, 393, 576, 735, 876, 996]

# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(users))

fig, ax = plt.subplots(figsize=(10, 7))
# 所有门店第一种产品的销量，注意控制柱子的宽度，这里选择0.25
ax.plot(xticks, all_block,  color='green', marker='x', linewidth=2, label="Share 10 Files")

ax.plot(xticks, all_block2,  color='#DE6B58', marker='x', linewidth=2, label="Share 50 Files")

# 显示数值
for a, b in zip(xticks, all_block):
    ax.text(a, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=12) 
# 显示数值
for a, b in zip(xticks, all_block2):
    ax.text(a, b+1, '%.0f'%b, ha='center', va= 'bottom', fontsize=12)


ax.set_title("Blocks Generated of Share Files", fontsize=15)
ax.set_xlabel("The Number of User", fontsize=12)
ax.set_ylabel("The Number of Block", fontsize=12)
ax.legend()

# 最后调整x轴标签的位置
ax.set_xticks(xticks )
ax.set_xticklabels(users)

plt.savefig('tmp.svg',dpi=800, bbox_inches='tight') # 保存成PDF放大后不失真（默认保存在了当前文件夹下）
plt.show()

