# 实验 分享50个文件的落块数量 折线图

import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")


# 用户数量
users = ["2", "4", "6", "8", "10", "12"]

# 总落块量
all_block = [194, 393, 576, 754, 920, 1098]

# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(users))

fig, ax = plt.subplots(figsize=(10, 7))
# 所有门店第一种产品的销量，注意控制柱子的宽度，这里选择0.25
ax.plot(xticks, all_block,  color='#DE6B58', marker='x', linewidth=1.5, label="blocks")

# 显示数值
for a, b in zip(xticks, all_block):
    ax.text(a, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=10) 


ax.set_title("All Blocks Of Share 10 Files", fontsize=15)
ax.set_xlabel("The number of user")
ax.set_ylabel("The number of block")
ax.legend()

# 最后调整x轴标签的位置
ax.set_xticks(xticks )
ax.set_xticklabels(users)

plt.savefig('tmp.png',dpi=800, bbox_inches='tight') # 保存成PDF放大后不失真（默认保存在了当前文件夹下）
plt.show()

