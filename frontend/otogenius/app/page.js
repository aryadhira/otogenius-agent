"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Fuel, Joystick, Receipt } from "lucide-react";
import { useState } from "react";
import { getRecommendation } from "./recommendation";
import { toast } from "sonner";
const Home = () => {
  const [prompt, setPrompt] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState([]);

  const handleSubmit = async () => {
    setLoading(true);
    setResult([]);

    const res = await getRecommendation({ input: prompt });
    console.log(res);
    if (res.status == 200) {
      setLoading(false);
      if (res?.data != null && res?.data?.length > 0) {
         toast("Success", {
           description: "Successfull load your car recommendation",
         });
        setResult(res.data);
      } else {
        toast("Result Not Found", {
          description: "Please refine the prompt!!",
        });
      }
    }
  };
  return (
    <div>
      <div className="flex flex-col gap-5 justify-center items-center">
        <div className="flex flex-col gap-2 justify-center items-center">
          <div className="text-2xl font-bold pt-10">Otogenius</div>
          <div className="text-md font-light">AI Used Car Consultant</div>
        </div>

        <div className="flex flex-col gap-2 justify-center items-center">
          <Textarea
            className={
              "max-w-[350] min-w-[350] md:max-w-[500] md:min-w-[500]  min-h-[250] p-5"
            }
            placeholder="describe your used car requirement in natural languange your description can contains some key eg: brand, model, car production year, budget, transmission type"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
          />
          <Button
            onClick={handleSubmit}
            disabled={loading}
            className={"w-full"}
          >
            {loading ? "Finding..." : "Find My Car"}
          </Button>
        </div>
      </div>
      <div className="flex justify-center items-center">
        <div className="grid grid-cols-1 md:grid-cols-3 xl:grid-cols-5 gap-5 p-10 items-center">
          {result.map((item, idx) => (
            <div key={idx} className="max-w-[300] min-h-[350]">
              <Card>
                <CardHeader>
                  <CardTitle>
                    {item.brand} {item.model}
                  </CardTitle>
                  <CardDescription>{item.production_year}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-col gap-5 text-sm">
                    <img
                      src={item.image_url || "https://placehold.co/300x300"}
                      onError={(e) =>
                        (e.currentTarget.src = "https://placehold.co/300x300")
                      }
                      alt={`${item.brand} ${item.model}`}
                      className="w-full h-48 object-cover rounded"
                    />
                    <div className="flex flex-col gap-1 text-xs">
                      <div className="flex gap-2 items-center">
                        <Fuel /> {item.fuel}
                      </div>
                      <div className="flex gap-2 items-center">
                        <Joystick /> {item.transmission}
                      </div>
                      <div className="flex gap-2 items-center">
                        <Receipt /> Rp. {item.price}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          ))}
        </div>
      </div>

      <div className="flex justify-center items-center font-light italic p-5">
        Recommendations are AI-generated. Please verify details before
        purchasing
      </div>
    </div>
  );
};

export default Home;
