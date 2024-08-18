import cv2 
from utils import create_cap, returnCameraIndexes

def query_and_open_available_cameras():
    print(returnCameraIndexes())

    select_index = int(input("Select camera index: "))

    cap = create_cap(select_index)

    if not cap.isOpened():
        print("Error: Could not open camera.")
        return

    last_working_frame = None
    broke = False

    # Loop to continuously read frames from the camera
    while True:
        # Read a frame from the camera
        ret, frame = cap.read()

        # Check if the frame was successfully read
        if ret:
            last_working_frame = frame
            if broke:
                print("Success: Read frame")
            broke = False
        elif not broke:
            print("Error: Could not read frame.")
            broke = True

        if last_working_frame is None:
            print("Error: No last working frame")
            break

        # Display the frame in a window named "Camera"
        cv2.imshow("Camera", last_working_frame)

        # Wait for 1 millisecond for keypress 'q' to quit the loop
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    # Release the camera and close the window
    cap.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    query_and_open_available_cameras()
