import { Moon, Sun } from "lucide-react";
import { useTheme } from "@/components/ThemeProvider";

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
      aria-label={`Switch to ${theme === "light" ? "dark" : "light"} mode`}
      className="p-2 rounded-md text-text-secondary hover:text-text-primary dark:hover:text-text-primary-dark transition-colors"
    >
      {theme === "light" ? <Moon size={18} /> : <Sun size={18} />}
    </button>
  );
}
