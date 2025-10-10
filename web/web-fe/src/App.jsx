import { motion, useMotionValue, useTransform } from 'framer-motion';
import AnimatedCounter from './components/AnimatedCounter'; // Giữ nguyên AnimatedCounter

// --- DATA CŨ (Không thay đổi) ---
const featuresData = [
    { icon: '🚗', title: 'Đặt xe nhanh chóng', desc: 'Tìm và đặt xe trong vài giây với giao diện đơn giản, thân thiện.' },
    { icon: '💰', title: 'Tiết kiệm chi phí', desc: 'Chia sẻ chi phí với bạn bè, giảm tới 70% so với đi xe riêng.' },
    { icon: '🛡️', title: 'An toàn tuyệt đối', desc: 'Xác thực sinh viên UIT, đánh giá tài xế và hành khách minh bạch.' },
    { icon: '📍', title: 'Theo dõi thời gian thực', desc: 'Xem vị trí xe và thời gian đến dự kiến ngay trên bản đồ.' },
    { icon: '👥', title: 'Cộng đồng UIT', desc: 'Kết nối với sinh viên cùng trường, tạo nhóm đi chung định kỳ.' },
    { icon: '🎯', title: 'Lộ trình thông minh', desc: 'Tối ưu tuyến đường, tiết kiệm thời gian và nhiên liệu hiệu quả.' }
];
const statsData = [
    { number: '5000+', label: 'Sinh viên' },
    { number: '500+', label: 'Tài xế' },
    { number: '20K+', label: 'Chuyến đi' },
    { number: '4.8★', label: 'Đánh giá' }
];
const testimonialsData = [
    { text: '"UIT-GO giúp mình tiết kiệm được rất nhiều chi phí đi lại. Mỗi tháng tiết kiệm được gần 1 triệu so với trước. Ứng dụng dễ dùng, tìm xe nhanh lắm!"', author: '- Nguyễn Văn A, Khoa CNTT' },
    { text: '"Làm tài xế UIT-GO giúp mình có thêm thu nhập mà vẫn đi đúng lộ trình đến trường. Kiếm được tiền xăng và có thêm bạn bè mới."', author: '- Trần Thị B, Khoa KHMT' },
    { text: '"An toàn và tiện lợi! Mình thích nhất là tính năng tạo nhóm đi chung cố định, giờ đi học với nhóm bạn thân mỗi ngày."', author: '- Lê Văn C, Khoa MMT&TT' }
];


// --- CÁC BIẾN THỂ ANIMATION CƠ BẢN (REUSABLE) ---
const fadeInDown = {
    hidden: { opacity: 0, y: -20 },
    visible: { opacity: 1, y: 0, transition: { duration: 0.6, ease: [0.2, 0.8, 0.2, 1] } }
};

const fadeInUp = {
    hidden: { opacity: 0, y: 30 },
    visible: { opacity: 1, y: 0, transition: { duration: 0.8, ease: [0.2, 0.8, 0.2, 1] } }
};

const staggerContainer = {
    hidden: {},
    visible: { transition: { staggerChildren: 0.15 } }
};

const slideInLeft = {
    hidden: { opacity: 0, x: -50 },
    visible: { opacity: 1, x: 0, transition: { duration: 0.8, ease: [0.2, 0.8, 0.2, 1] } }
};

const slideInRight = {
    hidden: { opacity: 0, x: 50 },
    visible: { opacity: 1, x: 0, transition: { duration: 0.8, ease: [0.2, 0.8, 0.2, 1] } }
};


function App() {
    // Giữ nguyên hiệu ứng 3D Parallax cho điện thoại
    const mouseX = useMotionValue(0);
    const mouseY = useMotionValue(0);
    // Điều chỉnh độ nhạy của hiệu ứng 3D
    const rotateX = useTransform(mouseY, [-400, 400], [15, -15]);
    const rotateY = useTransform(mouseX, [-400, 400], [-15, 15]);

    const handleMouseMove = (event) => {
        const { clientX, clientY, currentTarget } = event;
        const { left, top, width, height } = currentTarget.getBoundingClientRect();
        mouseX.set(clientX - left - width / 2);
        mouseY.set(clientY - top - height / 2);
    };

    const handleMouseLeave = () => {
        mouseX.set(0);
        mouseY.set(0);
    }

    return (
        <>
            {/* --- HERO SECTION --- */}
            <section className="hero" onMouseMove={handleMouseMove} onMouseLeave={handleMouseLeave}>
                <div className="hero-container">
                    {/* Hero Content với animation fadeIn */}
                    <motion.div className="hero-content" variants={staggerContainer} initial="hidden" animate="visible">
                        <motion.h1 variants={fadeInUp} className="logo">UIT-GO</motion.h1>
                        <motion.p variants={fadeInUp} className="tagline">
                            Đi chung thông minh, tiết kiệm chi phí<br />Dành riêng cho sinh viên UIT
                        </motion.p>
                        <motion.div variants={fadeInUp}>
                            <a href="#features" className="cta-button">
                                Khám phá ngay <span role="img" aria-label="arrow">→</span>
                            </a>
                        </motion.div>
                    </motion.div>
                    
                    {/* Phone Mockup với animation 3D và Floating */}
                    <motion.div className="phone-mockup-wrapper" style={{ rotateX, rotateY }} transition={{ type: 'spring', stiffness: 250, damping: 30 }}>
                        <motion.div
                            className="phone-mockup"
                            initial={{ opacity: 0, scale: 0.7 }}
                            animate={{ opacity: 1, scale: 1 }}
                            transition={{ duration: 1, ease: [0.2, 0.8, 0.2, 1], delay: 0.5 }}
                        >
                            <motion.div
                                animate={{ y: [0, -10, 0] }} // Floating nhẹ nhàng hơn
                                transition={{ duration: 5, ease: 'easeInOut', repeat: Infinity }}
                            >
                                <div className="screen">
                                    <div className="app-screen">
                                        <div className="header-logo"><span className="icon">🚘</span> UIT-GO</div>
                                        <div className="map-placeholder"><p>Bản đồ khu vực UIT</p><span className="indicator">📍</span></div>
                                        <div className="search-box">Tìm chuyến đi đến Thủ Đức...</div>
                                        <div className="action-buttons">
                                            <button className="action-btn offer-ride">Đăng chuyến</button>
                                            <button className="action-btn find-ride">Tìm chuyến</button>
                                        </div>
                                    </div>
                                </div>
                            </motion.div>
                        </motion.div>
                    </motion.div>
                </div>
            </section>
            
            {/* --- FEATURES SECTION - Bố cục xen kẽ hình ảnh và nội dung --- */}
            <section id="features" className="section">
                <div className="container">
                    <motion.h2 
                        className="section-title" 
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Tính năng nổi bật
                    </motion.h2>
                    <motion.p 
                        className="section-subtitle"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Trải nghiệm đi chung hiện đại, an toàn và tiện lợi
                    </motion.p>
                    
                    {/* Bố cục grid 2 cột cho feature cards */}
                    <motion.div 
                        className="feature-grid" 
                        variants={staggerContainer} 
                        initial="hidden" 
                        whileInView="visible" 
                        viewport={{ once: true, amount: 0.2 }}
                    >
                        {featuresData.map((feature, index) => (
                            <motion.div 
                                key={index} 
                                className="feature-card" 
                                variants={fadeInUp}
                                whileHover={{ y: -8, boxShadow: 'var(--shadow-heavy)' }}
                                transition={{ duration: 0.3 }}
                            >
                                <span className="feature-icon">{feature.icon}</span>
                                <h3 className="feature-title">{feature.title}</h3>
                                <p className="feature-desc">{feature.desc}</p>
                            </motion.div>
                        ))}
                    </motion.div>
                </div>
            </section>
            
            {/* --- STATS SECTION --- */}
            <section className="section">
                <div className="container">
                    <motion.h2 
                        className="section-title"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Con số ấn tượng
                    </motion.h2>
                    <motion.p 
                        className="section-subtitle"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        UIT-GO đang được tin dùng bởi cộng đồng sinh viên
                    </motion.p>
                    <motion.div 
                        className="stats-grid" 
                        variants={staggerContainer} 
                        initial="hidden" 
                        whileInView="visible" 
                        viewport={{ once: true, amount: 0.4 }}
                    >
                        {statsData.map((stat, index) => (
                            <motion.div key={index} className="stat-card" variants={fadeInUp}>
                                <div className="stat-number"><AnimatedCounter targetText={stat.number} /></div>
                                <div className="stat-label">{stat.label}</div>
                            </motion.div>
                        ))}
                    </motion.div>
                </div>
            </section>

            {/* --- TESTIMONIALS SECTION --- */}
            <section className="section">
                <div className="container">
                    <motion.h2 
                        className="section-title"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Sinh viên nói gì về UIT-GO
                    </motion.h2>
                    <motion.p 
                        className="section-subtitle"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Trải nghiệm thực tế từ cộng đồng
                    </motion.p>
                    <motion.div 
                        className="testimonials-container" 
                        variants={staggerContainer} 
                        initial="hidden" 
                        whileInView="visible" 
                        viewport={{ once: true, amount: 0.2 }}
                    >
                        {testimonialsData.map((testimonial, index) => (
                            <motion.div key={index} className="testimonial" variants={fadeInUp}>
                                <p className="testimonial-text">{testimonial.text}</p>
                                <p className="testimonial-author">{testimonial.author}</p>
                            </motion.div>
                        ))}
                    </motion.div>
                </div>
            </section>

            {/* --- CTA SECTION --- */}
            <section className="section cta-section">
                <div className="container">
                    <motion.h2 
                        className="section-title"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Sẵn sàng bắt đầu?
                    </motion.h2>
                    <motion.p 
                        className="section-subtitle"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Tải ứng dụng ngay hôm nay và nhận ưu đãi chuyến đi đầu tiên
                    </motion.p>
                    <motion.div
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInUp} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        <a href="#" className="cta-button">
                            Tải UIT-GO miễn phí <span role="img" aria-label="download">↓</span>
                        </a>
                    </motion.div>
                </div>
            </section>

            {/* --- FOOTER --- */}
            <footer className="footer">
                <div className="footer-content">
                    <motion.div 
                        className="footer-logo"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInDown} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        UIT-GO
                    </motion.div>
                    <motion.p 
                        className="footer-text"
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInDown} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        Giải pháp đi chung thông minh dành cho sinh viên UIT<br />An toàn • Tiết kiệm • Tiện lợi
                    </motion.p>
                    <motion.div 
                        className="social-links"
                        variants={staggerContainer} 
                        initial="hidden" 
                        whileInView="visible" 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        <motion.a href="#" className="social-link" variants={fadeInDown}>f</motion.a>
                        <motion.a href="#" className="social-link" variants={fadeInDown}>📷</motion.a>
                        <motion.a href="#" className="social-link" variants={fadeInDown}>💬</motion.a>
                    </motion.div>
                    <motion.p 
                        style={{ color: 'rgba(255,255,255,0.4)', marginTop: '40px' }}
                        initial="hidden" 
                        whileInView="visible" 
                        variants={fadeInDown} 
                        viewport={{ once: true, amount: 0.5 }}
                    >
                        © 2025 UIT-GO. All rights reserved.
                    </motion.p>
                </div>
            </footer>
        </>
    );
}

export default App;