# 实验 分享10个文件的落块数量 

import numpy as np
import matplotlib.pyplot as plt

# %matplotlib inline

# plt.style.use("ggplot")


# 用户数量
users = ["2", "4", "6", "8", "10", "12"]

# 初始化落块
init = [9, 11, 13, 15, 17, 19]
# 文件上传落块
upload = [20, 40, 58, 75, 91, 105]
# 文件分享落块
share = [20, 39, 57, 78, 90, 103]

# 创建分组柱状图，需要自己控制x轴坐标
xticks = np.arange(len(users))

fig, ax = plt.subplots(figsize=(10, 7))
# 所有门店第一种产品的销量，注意控制柱子的宽度，这里选择0.25
ax.bar(xticks, init, width=0.25, label="System Init", color="y")
# 所有门店第二种产品的销量，通过微调x轴坐标来调整新增柱子的位置
ax.bar(xticks + 0.25, upload, width=0.25, label="File Upload", color="green")
# 所有门店第三种产品的销量，继续微调x轴坐标调整新增柱子的位置
ax.bar(xticks + 0.5, share, width=0.25, label="File Share", color="red")

# 显示数值
for a, b in zip(xticks, init):
    ax.text(a, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=12) 
for a, b in zip(xticks, upload):
    ax.text(a+0.25, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=12) 
for a, b in zip(xticks, share):
    ax.text(a+0.5, b+1, '%.0f'%b, ha='center', va= 'bottom',fontsize=12) 

ax.set_title("Blocks Generated of Share 10 Files", fontsize=15)
ax.set_xlabel("The Number of User", fontsize=12)
ax.set_ylabel("The Number of Block", fontsize=12)
ax.legend()

# 最后调整x轴标签的位置
ax.set_xticks(xticks + 0.25)
ax.set_xticklabels(users)

plt.savefig('tmp.svg',dpi=800, bbox_inches='tight') # 保存成PDF放大后不失真（默认保存在了当前文件夹下）
plt.show()

