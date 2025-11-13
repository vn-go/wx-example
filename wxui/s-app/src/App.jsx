import { createSignal } from 'solid-js';
import './App.css';

function App() {
  const [user, setUser] = createSignal({
    name: "hello",
    email: "test@example.com"
  });
  return (
    <>
      <div>
        <button
          onClick={() => setUser({ name: "Updated!", email: "new@test.com" })}
          class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        >
          Update User
        </button>
      </div>
      <div>
        <pre class="mt-4 p-4 bg-gray-100 rounded text-xs">
          {JSON.stringify(user(), null, 2)}
        </pre>
      </div>
    </>
  )
}

export default App
