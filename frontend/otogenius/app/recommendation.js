"use server"

export const getRecommendation = async(param) =>{
    let result = []
    try {
      
      const res = await fetch("http://localhost:9009/api/recommendation", {
        method: "POST",
        body: JSON.stringify(param),
      });

      if (!res.ok) {
        throw new Error(`Error: ${res.status}`);
      }

      result = await res.json();
      
    } catch (e) {
      console.error("Failed to get recommendation:", e);
    }

    return result
}