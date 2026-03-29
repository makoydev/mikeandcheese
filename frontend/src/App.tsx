import { Routes, Route } from "react-router-dom";
import { Layout } from "@/components/Layout";
import { HomePage } from "@/pages/HomePage";
import { PostPage } from "@/pages/PostPage";
import { AboutPage } from "@/pages/AboutPage";
import { EditorPage } from "@/pages/EditorPage";

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/new" element={<EditorPage />} />
        <Route path="/edit/:slug" element={<EditorPage />} />
        <Route path="/post/:slug" element={<PostPage />} />
        <Route path="/about" element={<AboutPage />} />
      </Route>
    </Routes>
  );
}
