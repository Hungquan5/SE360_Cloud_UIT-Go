import { useEffect, useRef } from 'react';
import { animate, useInView } from 'framer-motion';

const AnimatedCounter = ({ targetText }) => {
    const ref = useRef(null);
    const isInView = useInView(ref, { once: true, margin: "-100px" });

    const isK = targetText.includes('K');
    const isPlus = targetText.includes('+');
    const isStar = targetText.includes('★');

    const numericalTarget = parseFloat(targetText.replace(/[K+★]/g, ''));

    useEffect(() => {
        if (isInView) {
            const controls = animate(0, numericalTarget, {
                duration: 2,
                onUpdate(value) {
                    let displayValue;
                    if (isStar) {
                        displayValue = value.toFixed(1) + '★';
                    } else if (isK) {
                        // Simplified logic for K+
                        displayValue = Math.floor(value).toLocaleString() + '+';
                    } else {
                        displayValue = Math.floor(value).toLocaleString() + (isPlus ? '+' : '');
                    }
                    ref.current.textContent = displayValue;
                },
                onComplete() {
                    // Ensure final value is exact
                    ref.current.textContent = targetText;
                }
            });
            return () => controls.stop();
        }
    }, [isInView, numericalTarget, targetText, isK, isPlus, isStar]);

    return <span ref={ref}>{isStar ? '0.0★' : '0'}</span>;
};

export default AnimatedCounter;