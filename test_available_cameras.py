import cv2 
import windows_capture_devices
from utils import create_cap

def returnCameraIndexes():
    devices = {}
    for i, d in enumerate(windows_capture_devices.get_capture_devices()):
        devices[d] = i

    return devices

def query_and_open_available_cameras():
    print(returnCameraIndexes())

    select_index = int(input("Select camera index: "))

    cap = create_cap(select_index)

    if not cap.isOpened():
        print("Error: Could not open camera.")
        return

    # Loop to continuously read frames from the camera
    while True:
        # Read a frame from the camera
        ret, frame = cap.read()

        # Check if the frame was successfully read
        if not ret:
            print("Error: Could not read frame.")
            break

        # Display the frame in a window named 'Webcam'
        cv2.imshow('Webcam', frame)

        # Wait for 1 millisecond for keypress 'q' to quit the loop
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    # Release the camera and close the window
    cap.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    query_and_open_available_cameras()
