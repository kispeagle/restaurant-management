import matplotlib.pyplot as plt
import numpy as np

def draw_sinx(id):
    sum = 0
    for c in id:
        sum += int(c)

    x = np.arange(0, sum, 0.1)   
    y = np.sin(x)

    plt.plot(x,y)
    plt.show()  
    
def draw_cosx(id):
    sum = 0
    for c in id:
        sum += int(c)

    x = np.arange(0, sum, 0.1)   

    y = np.cos(x)

    plt.plot(x,y)
    plt.show()  


def draw_tanx(id):
    sum = 0
    for c in id:
        sum += int(c)

    x = np.arange(0, sum, 0.1)   
    y = np.tan(x)

    plt.plot(x,y)
    plt.show()  



if __name__ == "__main__":
    import cv2

    img = cv2.imread("kimcuong.jpg")
    cv2.putText(img, "110570065", (50, 50), cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 0, 0), 2)
    cv2.rectangle(img, (0,0), (600,652), (255, 255, 0), 3)
    cv2.imshow("img", img)
    cv2.imwrite("kimcuong_result.jpg", img)
    cv2.waitKey(0)