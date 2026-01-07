'use client';

import { LoginForm } from '@/components/login-form';

export default function LoginPage() {
    return (
        <div className="container relative h-screen flex-col items-center justify-center md:grid md:max-w-none md:grid-cols-2 lg:px-0">
            <div className="relative hidden h-full flex-col bg-muted p-10 text-white lg:flex dark:border-r">
                <div className="absolute inset-0 bg-zinc-900" />
                <div
                    className="absolute inset-0 z-0 bg-zinc-900"
                    style={{
                        backgroundImage: 'radial-gradient(circle at center, #1e1b4b 0%, #000000 100%)',
                    }}
                />
                <div className="relative z-20 flex items-center text-lg font-medium">
                    <div className="w-8 h-8 rounded-lg bg-brand-600 mr-2 flex items-center justify-center text-white font-bold">
                        CP
                    </div>
                    Club Pulse
                </div>
                <div className="relative z-20 mt-auto">
                    <blockquote className="space-y-2">
                        <p className="text-lg">
                            &ldquo;This platform has completely transformed how we manage our club facilities. The premium scheduling features are a game changer.&rdquo;
                        </p>
                        <footer className="text-sm">Sofia Davis</footer>
                    </blockquote>
                </div>
            </div>
            <div className="lg:p-8 flex items-center justify-center h-full">
                <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
                    <LoginForm />
                </div>
            </div>
        </div>
    );
}
