import Footer from "@/components/layout/Footer";
import Header from "@/components/layout/Header";
import { TodoInput, TodoList } from "@/features/todo";

function App() {
  return (
    <div className="max-w-3xl mx-auto p-6 bg-white mt-12">
      <Header />
      <main>
        <TodoInput />
        <TodoList />
      </main>
      <Footer />
    </div>
  );
}

export default App;
