//  Spinner3 - три подскакивающие точки

export default function Spinner3() {
  const style = `
    @keyframes bounce {
      0%, 80%, 100% {
        opacity: 0.3;
        transform: scale(0.8);
      }
      40% {
        opacity: 1;
        transform: scale(1.2);
      }
    }
  `;
  
  return (
    <>
      <style>{style}</style>
      <div style={{
        display: 'flex',
        gap: '6px',
        justifyContent: 'center',
        alignItems: 'center'
      }}>
        <div style={{
          width: '12px',
          height: '12px',
          backgroundColor: 'brown',
          borderRadius: '50%',
          animation: 'bounce 1.4s infinite ease-in-out both',
          animationDelay: '-0.32s'
        }} />
        <div style={{
          width: '12px',
          height: '12px',
          backgroundColor: 'brown',
          borderRadius: '50%',
          animation: 'bounce 1.4s infinite ease-in-out both',
          animationDelay: '-0.16s'
        }} />
        <div style={{
          width: '12px',
          height: '12px',
          backgroundColor: 'brown',
          borderRadius: '50%',
          animation: 'bounce 1.4s infinite ease-in-out both'
        }} />
      </div>
    </>
  );
}
