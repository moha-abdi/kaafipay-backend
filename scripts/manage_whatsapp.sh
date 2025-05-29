#!/bin/bash

API_URL="http://localhost:8080/api/v1/admin/whatsapp"
ADMIN_TOKEN="bea305c8e8bcc770ca559940929a1047860316bdeac52d5bf90576d2f3277f74"  # Replace with your actual admin token
QR_IMAGE_PATH="$HOME/Pictures/whatsapp_qr.png"

function list_sessions() {
    curl -s -H "X-Admin-Token: $ADMIN_TOKEN" "$API_URL/sessions"
}

function get_session() {
    if [ -z "$1" ]; then
        echo "Usage: $0 get-session <session-id>"
        exit 1
    fi
    curl -s -H "X-Admin-Token: $ADMIN_TOKEN" "$API_URL/sessions/$1"
}

function display_qr() {
    local qr_data="$1"
    # Remove any "data:image/png;base64," prefix if present
    qr_data="${qr_data#*,}"
    # Remove any quotes and whitespace
    qr_data=$(echo "$qr_data" | tr -d '"' | tr -d '[:space:]')
    
    # Ensure Pictures directory exists
    mkdir -p "$HOME/Pictures"
    
    # Save QR code to Pictures directory
    echo "$qr_data" | base64 -d > "$QR_IMAGE_PATH"
    echo "QR code saved to: $QR_IMAGE_PATH"
    echo "Please open the image and scan it with WhatsApp"
}

function add_session() {
    if [ -z "$1" ]; then
        echo "Usage: $0 add-session <session-id>"
        exit 1
    fi

    echo "Creating new WhatsApp session: $1"
    echo "This will generate a QR code that you need to scan with WhatsApp..."
    
    response=$(curl -s -X POST \
        -H "X-Admin-Token: $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"session_id\": \"$1\", \"read_incoming_messages\": true, \"sync_full_history\": false}" \
        "$API_URL/sessions")
    
    # Extract QR code if present
    qr_code=$(echo "$response" | jq -r '.qr // empty')
    
    if [ ! -z "$qr_code" ]; then
        echo "QR code received. Saving to Pictures directory..."
        display_qr "$qr_code"
        
        # Start polling for session status
        echo "Waiting for QR code to be scanned..."
        while true; do
            sleep 5
            status_response=$(get_session "$1")
            session_status=$(echo "$status_response" | jq -r '.status // empty')
            
            if [ "$session_status" = "CONNECTED" ]; then
                echo "Session successfully connected!"
                rm -f "$QR_IMAGE_PATH"
                break
            elif [ "$session_status" = "FAILED" ]; then
                echo "Session creation failed!"
                rm -f "$QR_IMAGE_PATH"
                exit 1
            fi
            echo "Waiting for connection... (status: $session_status)"
        done
    else
        error_msg=$(echo "$response" | jq -r '.error // empty')
        if [ ! -z "$error_msg" ]; then
            echo "Error: $error_msg"
            exit 1
        fi
        echo "Session created successfully!"
    fi
}

function delete_session() {
    if [ -z "$1" ]; then
        echo "Usage: $0 delete-session <session-id>"
        exit 1
    fi
    curl -s -X DELETE -H "X-Admin-Token: $ADMIN_TOKEN" "$API_URL/sessions/$1"
}

# Check for required commands
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Please install jq first."
    exit 1
fi

case "$1" in
    "list")
        list_sessions
        ;;
    "get")
        get_session "$2"
        ;;
    "add")
        add_session "$2"
        ;;
    "delete")
        delete_session "$2"
        ;;
    *)
        echo "Usage: $0 {list|get|add|delete} [session-id]"
        exit 1
        ;;
esac 