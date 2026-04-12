import { w as writable } from "./index.js";
function createToasts() {
  const { subscribe, update } = writable([]);
  function add(message, type = "info") {
    const id = crypto.randomUUID();
    update((list) => [...list, { id, message, type }]);
    return id;
  }
  function remove(id) {
    update((list) => list.filter((t) => t.id !== id));
  }
  return {
    subscribe,
    success: (msg) => add(msg, "success"),
    error: (msg) => add(msg, "error"),
    info: (msg) => add(msg, "info"),
    remove
  };
}
const toasts = createToasts();
export {
  toasts as t
};
