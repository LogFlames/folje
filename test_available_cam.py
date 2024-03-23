import cv2 

def returnCameraIndexes():
    index = 0
    arr = []
    i = 10
    while i > 0:
        cap = cv2.VideoCapture(index)
        if cap.read()[0]:
            arr.append(index)
            cap.release()
        index += 1
        i -= 1
    return arr

def main():
    print(returnCameraIndexes())

    select_index = int(input("Select camera index: "))

    cap = cv2.VideoCapture(select_index)
    cv2.namedWindow("Webcam")

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
    main()

