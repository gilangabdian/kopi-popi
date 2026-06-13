import { Coffee } from "lucide-react";
export default function Footer() {
  return (
    <div>
      {/* Footer */}
      <footer className="bg-foreground text-background py-10 text-center">
        <div className="flex items-center justify-center gap-2 mb-4">
          <Coffee className="w-6 h-6 text-primary" />
          <span className="text-xl font-bold">Kopi Popi</span>
        </div>
        <p className="text-muted/60 text-sm">© {new Date().getFullYear()} Kopi Popi. All rights reserved.</p>
      </footer>
    </div>
  );
}
