import { ref } from "vue";

export interface NotificationMessage {
  id: string;
  userMessage: string;
  type: "error" | "success";
  technicalDetails?: string;
  timestamp: Date;
  autoHide?: boolean;
  autoHideDelay?: number;
}

// Global reactive state for notification management
const currentNotification = ref<NotificationMessage | null>(null);
const notificationHistory = ref<NotificationMessage[]>([]);

export function useNotificationManager() {
  /**
   * Show an error to the user with a user-friendly message
   * Technical details are logged to console for debugging
   * Errors never auto-hide and must be manually dismissed
   */
  function showError(userMessage: string, technicalDetails?: unknown) {
    // Log technical details to console for developers
    if (technicalDetails) {
      console.error("Error occurred:", {
        userMessage,
        technicalDetails:
          technicalDetails instanceof Error
            ? technicalDetails.message
            : String(technicalDetails),
        stack:
          technicalDetails instanceof Error
            ? technicalDetails.stack
            : undefined,
        timestamp: new Date(),
      });
    } else {
      console.error("Error occurred:", userMessage);
    }

    // Create error message for user display
    const errorMessage: NotificationMessage = {
      id: Date.now().toString(),
      userMessage,
      type: "error",
      technicalDetails:
        technicalDetails instanceof Error
          ? technicalDetails.message
          : String(technicalDetails),
      timestamp: new Date(),
      autoHide: false, // Errors never auto-hide
    };

    // Update global state
    currentNotification.value = errorMessage;
    notificationHistory.value.push(errorMessage);
  }

  /**
   * Clear the current notification message
   */
  function clearError() {
    currentNotification.value = null;
  }

  /**
   * Show a success message
   * Success messages auto-hide by default after a configurable delay
   */
  function showSuccess(
    message: string,
    options: { autoHide?: boolean; autoHideDelay?: number } = {},
  ) {
    console.log("Success:", message);

    const { autoHide = true, autoHideDelay = 3000 } = options;

    const successMessage: NotificationMessage = {
      id: Date.now().toString(),
      userMessage: message,
      type: "success",
      timestamp: new Date(),
      autoHide,
      autoHideDelay,
    };

    // Update global state
    currentNotification.value = successMessage;
    notificationHistory.value.push(successMessage);

    // Auto-hide if configured
    if (autoHide) {
      setTimeout(() => {
        // Only clear if this is still the current notification
        if (currentNotification.value?.id === successMessage.id) {
          clearError();
        }
      }, autoHideDelay);
    }
  }

  return {
    currentNotification,
    showError,
    clearError,
    showSuccess,
  };
}
