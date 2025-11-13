import { signal } from "@preact/signals";

const count = signal(0);
const user = signal({ name: "hello", email: "test" });

export default function App() {
  return (
    <div class="p-8 bg-gray-50 min-h-screen">
      <p>Count: <strong>{count}</strong></p>
      <p>Name: <strong>{user.value.name}</strong></p>
      <button onClick={() => count.value++} class="m-1 px-3 py-1 bg-blue-500 text-white rounded">+1</button>
      <button onClick={() => user.value = { ...user.value, name: "123456" }} class="m-1 px-3 py-1 bg-green-500 text-white rounded">
        Update
      </button>
    </div>
  );
}