import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Coffee, CreditCard, Gift } from "lucide-react";

export default function LandingPage() {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Navbar */}
      <header className="px-6 py-4 border-b border-secondary/20 flex justify-between items-center bg-background sticky top-0 z-10">
        <div className="flex items-center gap-2">
          <Coffee className="w-8 h-8 text-primary" />
          <span className="text-xl font-bold text-foreground">Kopi Popi</span>
        </div>
        <div className="flex gap-4">
          <Link href="/login">
            <Button variant="outline" className="border-secondary text-secondary-foreground hover:bg-secondary hover:text-white">
              Login
            </Button>
          </Link>
          <Link href="/register">
            <Button className="bg-primary text-primary-foreground hover:bg-primary/90">
              Daftar
            </Button>
          </Link>
        </div>
      </header>

      {/* Hero Section */}
      <main className="flex-1 flex flex-col items-center justify-center text-center px-4 py-20 bg-gradient-to-b from-primary/10 to-background">
        <Badge className="mb-4 bg-secondary text-white hover:bg-secondary/80">Kini Hadir di Kotamu!</Badge>
        <h1 className="text-5xl md:text-7xl font-extrabold text-foreground tracking-tight max-w-4xl mb-6">
          Nikmati Kopi Terbaik <br />
          <span className="text-primary">Kapanpun, Di Manapun.</span>
        </h1>
        <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mb-10">
          Kopi Popi menghadirkan perpaduan biji kopi nusantara dengan teknologi pemesanan modern.
          Dapatkan poin loyalitas setiap transaksi!
        </p>
        <div className="flex gap-4">
          <Button size="lg" className="bg-primary text-primary-foreground hover:bg-primary/90 text-lg px-8">
            Pesan Sekarang
          </Button>
          <Button size="lg" variant="outline" className="text-lg px-8 border-secondary text-secondary-foreground hover:bg-secondary/10">
            Pelajari Lebih Lanjut
          </Button>
        </div>
      </main>

      {/* Features Section */}
      <section className="py-20 px-6 bg-muted/20">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-foreground mb-4">Kenapa Memilih Kopi Popi?</h2>
            <p className="text-muted-foreground">Kami tidak hanya menjual kopi, kami memberikan pengalaman.</p>
          </div>
          
          <div className="grid md:grid-cols-3 gap-8">
            <Card className="border-secondary/20 shadow-sm hover:shadow-md transition-shadow">
              <CardHeader className="text-center pb-2">
                <div className="w-16 h-16 rounded-full bg-primary/20 flex items-center justify-center mx-auto mb-4">
                  <Coffee className="w-8 h-8 text-primary" />
                </div>
                <CardTitle className="text-xl">Kopi Pilihan</CardTitle>
              </CardHeader>
              <CardContent className="text-center text-muted-foreground">
                <p>Biji kopi nusantara yang disangrai dengan sempurna untuk menghasilkan rasa yang kaya dan otentik.</p>
              </CardContent>
            </Card>

            <Card className="border-secondary/20 shadow-sm hover:shadow-md transition-shadow">
              <CardHeader className="text-center pb-2">
                <div className="w-16 h-16 rounded-full bg-primary/20 flex items-center justify-center mx-auto mb-4">
                  <Gift className="w-8 h-8 text-primary" />
                </div>
                <CardTitle className="text-xl">Poin Loyalitas</CardTitle>
              </CardHeader>
              <CardContent className="text-center text-muted-foreground">
                <p>Kumpulkan poin di setiap transaksi dan tukarkan dengan kopi gratis atau potongan harga menarik.</p>
              </CardContent>
            </Card>

            <Card className="border-secondary/20 shadow-sm hover:shadow-md transition-shadow">
              <CardHeader className="text-center pb-2">
                <div className="w-16 h-16 rounded-full bg-primary/20 flex items-center justify-center mx-auto mb-4">
                  <CreditCard className="w-8 h-8 text-primary" />
                </div>
                <CardTitle className="text-xl">Bayar Mudah</CardTitle>
              </CardHeader>
              <CardContent className="text-center text-muted-foreground">
                <p>Terima berbagai metode pembayaran mulai dari tunai, QRIS, hingga dompet digital.</p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-foreground text-background py-10 text-center">
        <div className="flex items-center justify-center gap-2 mb-4">
          <Coffee className="w-6 h-6 text-primary" />
          <span className="text-xl font-bold">Kopi Popi</span>
        </div>
        <p className="text-muted/60 text-sm">
          © {new Date().getFullYear()} Kopi Popi. All rights reserved.
        </p>
      </footer>
    </div>
  );
}
