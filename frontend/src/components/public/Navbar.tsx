import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Coffee } from "lucide-react";

export default function Navbar() {
  return (
    <div>
      {/* Navbar */}
      <header className="px-6 py-4 border-b border-secondary/20 flex justify-between items-center bg-background sticky top-0 z-10">
        <div className="flex items-center gap-2">
          <Coffee className="w-8 h-8 text-primary" />
          <span className="text-xl font-bold text-foreground">Kopi Popi</span>
        </div>
        <div className="flex gap-4">
          <Link href="/login">
            <Button
              variant="outline"
              className="border-secondary text-secondary-foreground hover:bg-secondary hover:text-white">
              Login
            </Button>
          </Link>
          <Link href="/register">
            <Button className="bg-primary text-primary-foreground hover:bg-primary/90">Daftar</Button>
          </Link>
        </div>
      </header>
    </div>
  );
}
