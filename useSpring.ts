/**
 * Spring animation presets untuk iOS-like natural motion
 * menggunakan motion-v library dengan physics-based animations
 */

export const useSpringAnimations = () => {
    /**
     * Spring preset untuk button press dengan bounce effect
     */
    const buttonPress = {
        initial: { scale: 1 },
        tap: { scale: 0.97 },
        transition: {
            type: 'spring',
            stiffness: 300,
            damping: 25,
        },
    };

    /**
     * Spring preset untuk card entrance dengan stagger
     */
    const cardEntrance = {
        initial: { opacity: 0, y: 20, scale: 0.95 },
        enter: { opacity: 1, y: 0, scale: 1 },
        transition: {
            type: 'spring',
            stiffness: 300,
            damping: 25,
            mass: 0.5,
        },
    };

    /**
     * Spring preset untuk fade in dengan natural movement
     */
    const fadeIn = {
        initial: { opacity: 0 },
        enter: { opacity: 1 },
        transition: {
            type: 'spring',
            stiffness: 200,
            damping: 20,
        },
    };

    /**
     * Spring preset untuk slide in dari kiri
     */
    const slideInLeft = {
        initial: { x: -100, opacity: 0 },
        enter: { x: 0, opacity: 1 },
        leave: { x: -100, opacity: 0 },
        transition: {
            type: 'spring',
            stiffness: 300,
            damping: 30,
        },
    };

    /**
     * Spring preset untuk slide in dari kanan
     */
    const slideInRight = {
        initial: { x: 100, opacity: 0 },
        enter: { x: 0, opacity: 1 },
        leave: { x: 100, opacity: 0 },
        transition: {
            type: 'spring',
            stiffness: 300,
            damping: 30,
        },
    };

    /**
     * Spring preset untuk modal/dialog dengan scale
     */
    const modal = {
        initial: { scale: 0.95, opacity: 0, y: 20 },
        enter: { scale: 1, opacity: 1, y: 0 },
        leave: { scale: 0.95, opacity: 0, y: 20 },
        transition: {
            type: 'spring',
            stiffness: 300,
            damping: 30,
        },
    };

    /**
     * Generate stagger delay untuk list items
     */
    const getStaggerDelay = (index: number, baseDelay: number = 50) => {
        return index * baseDelay;
    };

    return {
        buttonPress,
        cardEntrance,
        fadeIn,
        slideInLeft,
        slideInRight,
        modal,
        getStaggerDelay,
    };
};
