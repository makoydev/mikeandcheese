import { Link, Outlet } from "react-router-dom";
import { ThemeToggle } from "@/components/ThemeToggle";

export function Layout() {
  return (
    <div className="min-h-screen flex flex-col bg-bg-light dark:bg-bg-dark">
      <header className="border-b border-border dark:border-border-dark">
        <div className="max-w-[720px] mx-auto px-6 py-5 flex items-center justify-between">
          <Link
            to="/"
            className="font-serif text-xl font-bold text-text-primary dark:text-text-primary-dark hover:opacity-80 transition-opacity"
          >
            mikeandcheese
          </Link>
          <nav className="flex items-center gap-6">
            <Link
              to="/"
              className="text-sm font-medium text-text-secondary hover:text-text-primary dark:hover:text-text-primary-dark transition-colors"
            >
              Home
            </Link>
            <Link
              to="/about"
              className="text-sm font-medium text-text-secondary hover:text-text-primary dark:hover:text-text-primary-dark transition-colors"
            >
              About
            </Link>
            <ThemeToggle />
          </nav>
        </div>
      </header>

      <main className="flex-1 max-w-[720px] w-full mx-auto px-6 py-10">
        <Outlet />
      </main>

      <footer className="border-t border-border dark:border-border-dark">
        <div className="max-w-[720px] mx-auto px-6 py-6 text-center text-sm text-text-secondary">
          &copy; {new Date().getFullYear()} mikeandcheese. All rights reserved.
        </div>
      </footer>
    </div>
  );
}
