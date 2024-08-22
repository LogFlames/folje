import cv2

from utils import create_cap, returnCameraIndexes


def query_and_open_available_cameras():
    print(returnCameraIndexes())

    select_index = int(input("Select camera index: "))

    cap = create_cap(select_index)

    if not cap.isOpened():
        print("Error: Could not open camera.")
        return

    while True:
        ret, frame = cap.read()

        if not ret:
            print("Error: Could not read frame.")
            break

        cv2.imshow("Camera", frame)

        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    # Release the camera and close the window
    cap.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    query_and_open_available_cameras()
