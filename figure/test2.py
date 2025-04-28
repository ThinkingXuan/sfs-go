import numpy as np
from scipy.interpolate import make_interp_spline
import matplotlib.pyplot as plt 

x=np.array([1,2,3,4,5,6,7])
y=np.array([100,50,25,12.5,6.25,3.125,1.5625])

model=make_interp_spline(x, y)

xs=np.linspace(1,7,500)
ys=model(xs)

plt.plot(xs, ys)
plt.title("Smooth Spline Curve")
plt.xlabel("X")
plt.ylabel("Y")
plt.show()