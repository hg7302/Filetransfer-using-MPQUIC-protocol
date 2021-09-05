import time
import pathlib

from matplotlib.animation import FuncAnimation
import matplotlib.pyplot as plt
import matplotlib.image as img


frame_counter = 0

start_time = time.time()
current_time = time.time()


def grab_frame():
    global frame_counter

    file = pathlib.Path('sample/img' + str(frame_counter) + '.jpg')

    missing_frame_counter = 0
    while (not file.exists() and missing_frame_counter <= 100):

        # check for the frame every TIMEOUT seconds
        TIMEOUT = 0.2
        time.sleep(TIMEOUT)
        missing_frame_counter += 1

    image = img.imread('sample/img' + str(frame_counter) + '.jpg')
    frame_counter += 1

    current_time = time.time()
    print("fps: ", frame_counter / (current_time - start_time))

    # delete the frame from the disk.
    file.unlink()
    return image

# create axes
ax1 = plt.subplot(111)

# create axes
im1 = ax1.imshow(grab_frame())

def update(i):
    im1.set_data(grab_frame())

# Animate the view every 25 miliseconds
ani = FuncAnimation(plt.gcf(), update, interval=25)
plt.show()
