import { ref } from "vue";

type Level = "info" | "warning" | "error";
interface Options {
  logMessage?: unknown;
  autoHide: boolean;
  autoHideDelay: number;
}
const defaultOptions: Options = { autoHide: true, autoHideDelay: 3000 };

export interface NotificationMessage {
  id: string;
  level: Level;
  userMessage: string;
}

// Global reactive state for notification management
const notifications = ref<NotificationMessage[]>([]);

export function useNotificationManager() {
  function show(
    level: Level,
    userMessage: string,
    options: Partial<Options> = {},
  ): NotificationMessage {
    const finalOptions: Options = { ...defaultOptions, ...options };

    // Log message
    if (finalOptions.logMessage !== undefined) {
      console.log(finalOptions.logMessage);
    }

    const notification = <NotificationMessage>{
      id: Date.now().toString(),
      level: level,
      userMessage: userMessage,
    };

    notifications.value.push(notification);

    if (finalOptions.autoHide) {
      setTimeout(() => clear(notification.id), finalOptions.autoHideDelay);
    }
    return notification;
  }

  function clear(id?: string) {
    if (id === undefined) {
      notifications.value = [];
      return;
    }
    const index = notifications.value.findIndex(
      (n: NotificationMessage) => n.id === id,
    );

    // index not found
    if (index === -1) {
      return;
    }
    notifications.value.splice(index, 1);
  }

  return {
    notifications,
    show,
    clear,
  };
}
