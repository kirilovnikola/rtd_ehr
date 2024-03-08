import websocket
import _thread
import time
import rel
import sys

def on_message(ws, message):
    print(message)

def on_error(ws, error):
    print(error)

def on_close(ws, close_status_code, close_msg):
    print("Connection closed (subscription id not recognized)")
    sys.exit()

def on_open(ws):
    print("Opened connection")

if __name__ == "__main__":
    ws = websocket.WebSocketApp("ws://localhost:8080/websocket",
                              on_open=on_open,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)

    ws.run_forever(dispatcher=rel, reconnect=5)  # Set dispatcher to automatic reconnection, 5 second reconnect delay if connection closed unexpectedly
    try:
        ws.send("bind "+sys.argv[1]) #Send fhir bind message with Subscription id
    except:
        print("Problem with bind subscription id")
        sys.exit()
    rel.signal(2, rel.abort)  # Keyboard Interrupt
    rel.dispatch()