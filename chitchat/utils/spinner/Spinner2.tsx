// Spinner2 - круговой прогресс-бар

export default function Spinner2() {
  const style = `
    @keyframes spin {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }
  `;
  
  return (
    <>
      <style>{style}</style>
      <div style={{
        width: '50px',
        height: '50px',
        border: '4px solid #f3f3f3',
        borderTop: '4px solid brown',
        borderRadius: '50%',
        animation: 'spin 1s linear infinite'
      }} />
    </>
  );
}
