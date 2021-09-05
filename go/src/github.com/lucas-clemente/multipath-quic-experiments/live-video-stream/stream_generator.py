import io
import time
import picamera
import socket


# This script generated a steady stream of live video frames from raspi camera module.
# Then it stores the the frame in-memory (bytes stream) and then finally send the frame
# to a local go client using TCP socket (using sockets for interprocess communication).


IP = "localhost"
PORT = 8002


sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
sock.bind((IP, PORT))
sock.listen(1)
print('Listening at', sock.getsockname())

sc, sockname = sock.accept()

frames = 0
start_time = time.time()


with picamera.PiCamera() as camera:

    # camera resolution
    camera.resolution = (300, 300)
    
    # creating a byte stream to store image
    stream = io.BytesIO()


    # Infinite loop which captures frames in 'jpeg' format.
    # use_video_port=True - This flag uses video port instead of capture port. It increases the frame rate,
    # but decreases the sharpness of the frame.
    for _ in camera.capture_continuous(stream, format='jpeg', use_video_port=True):

        # Truncate the stream to the current position (in case
        # prior iterations output a longer image)
        stream.truncate()
        stream.seek(0)
        print("frame size in bytes: ", len(stream.getbuffer()))

        # limiting FPS to sync with output stream
        time.sleep(0.05)

        # sending size of the frame first
        buffer_size = str(len(stream.getbuffer()))
        buffer_size_data = buffer_size.ljust(20, ':')
        sc.sendall(buffer_size_data.encode())

        # sending the actual frame
        sc.sendall(stream.getbuffer())

        frames += 1
        current_time = time.time()
        fps = frames / (current_time - start_time)
        print("fps: ", fps)
