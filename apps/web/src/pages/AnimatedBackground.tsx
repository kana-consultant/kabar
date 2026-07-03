import { useEffect, useRef } from "react";

export default function AnimatedBackground() {
    const canvasRef = useRef<HTMLCanvasElement>(null);

    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) return;

        const ctx = canvas.getContext('2d');
        if (!ctx) return;

        let width = canvas.width = canvas.parentElement?.clientWidth || window.innerWidth;
        let height = canvas.height = canvas.parentElement?.clientHeight || window.innerHeight;
        let animationFrameId: number;
        let time = 0;

        // Mouse tracking & Target positioning
        const mouse = { x: width / 2, y: height / 2, targetX: width / 2, targetY: height / 2, active: false };

        const resizeCanvas = () => {
            if (!canvas || !canvas.parentElement) return;
            width = canvas.width = canvas.parentElement.clientWidth;
            height = canvas.height = canvas.parentElement.clientHeight;
        };

        const handleMouseMove = (e: MouseEvent) => {
            const rect = canvas.getBoundingClientRect();
            mouse.targetX = e.clientX - rect.left;
            mouse.targetY = e.clientY - rect.top;
            mouse.active = true;
        };

        const handleMouseLeave = () => {
            mouse.active = false;
        };

        // Definisikan struktur partikel abstrak tunggal
        class AbstractParticle {
            x: number;
            y: number;
            baseX: number;
            baseY: number;
            size: number;
            speedX: number;
            speedY: number;
            angle: number;
            spin: number;
            hue: number;
            alpha: number;

            constructor() {
                this.x = Math.random() * width;
                this.y = Math.random() * height;
                this.baseX = this.x;
                this.baseY = this.y;
                this.size = Math.random() * 2 + 0.5;
                this.speedX = (Math.random() - 0.5) * 0.8;
                this.speedY = (Math.random() - 0.5) * 0.8;
                this.angle = Math.random() * Math.PI * 2;
                this.spin = (Math.random() - 0.5) * 0.01;
                this.hue = 240 + Math.random() * 50; // Range warna ungu ke indigo/violet
                this.alpha = Math.random() * 0.5 + 0.2;
            }

            update(time: number) {
                this.angle += this.spin;

                // Flow Field matematika dasar (membuat efek cairan/gerakan abstrak bergelombang)
                const flowX = Math.sin(this.angle + time * 0.002) * 0.5;
                const flowY = Math.cos(this.angle - time * 0.001) * 0.5;

                this.x += this.speedX + flowX;
                this.y += this.speedY + flowY;

                // Interaksi interaktif dengan Mouse (jika mouse aktif bergerak di area hero)
                if (mouse.active) {
                    const dx = mouse.x - this.x;
                    const dy = mouse.y - this.y;
                    const distance = Math.sqrt(dx * dx + dy * dy);

                    if (distance < 180) {
                        const force = (180 - distance) / 180;
                        // Partikel terseret perlahan mengikuti pusaran arus mouse
                        this.x += (dx / distance) * force * 1.5;
                        this.y += (dy / distance) * force * 1.5;
                    }
                }

                // Boundary reset jika partikel keluar dari area canvas
                if (this.x < 0 || this.x > width || this.y < 0 || this.y > height) {
                    this.x = Math.random() * width;
                    this.y = Math.random() * height;
                    this.angle = Math.random() * Math.PI * 2;
                }
            }

            draw(ctx: CanvasRenderingContext2D) {
                ctx.beginPath();
                ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
                ctx.fillStyle = `hsla(${this.hue}, 85%, 70%, ${this.alpha})`;
                ctx.fill();
            }
        }

        // Generate kumpulan partikel abstrak (jumlah 120 pas untuk estetika bersih & ringan)
        const particleCount = 120;
        const particles: AbstractParticle[] = [];
        for (let i = 0; i < particleCount; i++) {
            particles.push(new AbstractParticle());
        }

        const animate = () => {
            time += 1;

            // Efek trail pudar (ghost effect) untuk menciptakan kesan motion blur yang halus
            ctx.fillStyle = 'rgba(9, 5, 20, 0.08)';
            ctx.fillRect(0, 0, width, height);

            // Interpolasi posisi mouse agar pergerakan responsif tapi tetap halus
            mouse.x += (mouse.targetX - mouse.x) * 0.08;
            mouse.y += (mouse.targetY - mouse.y) * 0.08;

            // Render partikel & garis penghubung abstrak yang samar
            for (let i = 0; i < particles.length; i++) {
                particles[i].update(time);
                particles[i].draw(ctx);

                // Buat garis koneksi antar partikel terdekat secara acak (Constellation Mesh effect)
                for (let j = i + 1; j < particles.length; j++) {
                    const dx = particles[i].x - particles[j].x;
                    const dy = particles[i].y - particles[j].y;
                    const dist = Math.sqrt(dx * dx + dy * dy);

                    if (dist < 65) {
                        ctx.beginPath();
                        ctx.moveTo(particles[i].x, particles[i].y);
                        ctx.lineTo(particles[j].x, particles[j].y);
                        const edgeAlpha = (1 - dist / 65) * 0.07;
                        ctx.strokeStyle = `hsla(${particles[i].hue}, 80%, 75%, ${edgeAlpha})`;
                        ctx.lineWidth = 0.5;
                        ctx.stroke();
                    }
                }
            }

            animationFrameId = requestAnimationFrame(animate);
        };

        animate();

        window.addEventListener('resize', resizeCanvas);
        window.addEventListener('mousemove', handleMouseMove);
        window.addEventListener('mouseleave', handleMouseLeave);

        return () => {
            cancelAnimationFrame(animationFrameId);
            window.removeEventListener('resize', resizeCanvas);
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseleave', handleMouseLeave);
        };
    }, []);

    return (
        <canvas
            ref={canvasRef}
            className="absolute inset-0 w-full h-full mix-blend-screen"
            style={{ pointerEvents: 'none' }}
        />
    );
};