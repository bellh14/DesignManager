
import star.base.report.Report;
import star.common.*;
import star.meshing.*;
import star.surfacewrapper.SurfaceWrapperAutoMeshOperation;
import star.vis.Scene;

import java.io.BufferedWriter;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.concurrent.TimeUnit;

/*
    error codes:
    1: failed to read input csv
    2: meshing error
    3: results saving error
    4: other error (probably sim error from mesh)
 */


public class DMPareto extends StarMacro {
    
    Double biplane1AOA = 0.0;
    Double biplane2AOA = 0.0;
    Double biplaneGapSize = 0.0;
    Double biplanePosition = 0.0;
    Double fw4thAOA = 0.0;

    @Override
    public void execute() {
        Simulation sim = getActiveSimulation();
        String baseDir = sim.getSessionDir();
        String simName = sim.getPresentationName();

        try {

            ReadCSVInputs(resolvePath("sim_inputs.csv"));

            long startTotalTime = System.nanoTime(); // will measure the total time taken of the sim
            updateSimParameters(sim);

            if (!updateMesh(sim)) { // runs meshing pipeline, catches errors
                System.out.println("Fatal Mesh Error");
                continue;
            }

            long iterationStartTime = System.nanoTime();
            sim.getSimulationIterator().run();

            long iterationEndTime = System.nanoTime();
            long iterationElapsedTime = iterationEndTime - iterationStartTime;

            System.out.println("Iteration Time Take: "
                    + TimeUnit.MINUTES.convert((iterationElapsedTime), TimeUnit.NANOSECONDS));

            saveScenes(sim, baseDir, simName);
            long endTotal = System.nanoTime();
            long totalElapsed = endTotal - startTotalTime;
            System.out.println("Total Time Taken: " + TimeUnit.MINUTES.convert(totalElapsed, TimeUnit.NANOSECONDS));
        } catch (Exception e) {
            e.printStackTrace();
            System.out.println("It is broken but probably not my fault");
            saveScenes(sim, baseDir, simName);
        }

    }

    public void updateSimParameters(Simulation sim) {


        ScalarGlobalParameter Biplane1AOA = ((ScalarGlobalParameter) sim.get(GlobalParameterManager.class)
                .getObject("Biplane1 AOA"));
        Units Biplane1AOAUnis = ((Units) sim.getUnitsManager().getObject("deg"));
        Biplane1AOA.getQuantity().setValueAndUnits(biplane1AOA, Biplane1AOAUnis);

        ScalarGlobalParameter Biplane2AOA = ((ScalarGlobalParameter) sim.get(GlobalParameterManager.class)
                .getObject("Biplane2 AOA"));
        Units Biplane2AOAUnis = ((Units) sim.getUnitsManager().getObject("deg"));
        Biplane2AOA.getQuantity().setValueAndUnits(biplane2AOA, Biplane2AOAUnis);

        ScalarGlobalParameter BiplaneGapSize = ((ScalarGlobalParameter) sim.get(GlobalParameterManager.class)
                .getObject("BiplaneGapSize"));
        Units BiplaneGapSizeUnits = ((Units) sim.getUnitsManager().getObject("mm"));
        Biplane1AOA.getQuantity().setValueAndUnits(biplaneGapSize, BiplaneGapSizeUnits);

        ScalarGlobalParameter BiplanePosition = ((ScalarGlobalParameter) sim.get(GlobalParameterManager.class)
                .getObject("BiplanePosition"));
        Units BiplanePositionUnits = ((Units) sim.getUnitsManager().getObject("mm"));
        Biplane1AOA.getQuantity().setValueAndUnits(biplanePosition, BiplanePositionUnits);
        
        ScalarGlobalParameter FW4THAOA = ((ScalarGlobalParameter) sim.get(GlobalParameterManager.class)
                .getObject("FW 4th Elemenent AOA"));
        Units FW4THAOAUnis = ((Units) sim.getUnitsManager().getObject("deg"));
        FW4THAOA.getQuantity().setValueAndUnits(fw4thAOA, FW4THAOAUnis);

        System.out.println("Biplane1 AOA: " + biplane1AOA);
        System.out.println("Biplane2 AOA: " + biplane2AOA);
        System.out.println("Biplane Gap Size: " + biplaneGapSize);
        System.out.println("Biplane Position: " + biplanePosition);
        System.out.println("4th Wing AOA: " + fw4thAOA);
    }

    public boolean updateMesh(Simulation sim) {
        try {
            long meshStartTime = System.nanoTime();
            MeshPipelineController mesh = sim.get(MeshPipelineController.class);
            mesh.clearGeneratedMeshes();

            sim.get(MeshOperationManager.class).executeAll();

            long meshEndTime = System.nanoTime();
            long meshElapsedTime = meshEndTime - meshStartTime;
            System.out
                    .println("Mesh pipeline time: " + TimeUnit.MINUTES.convert(meshElapsedTime, TimeUnit.NANOSECONDS));
        } catch (Exception e) { // catches fatal mesh errors
            e.printStackTrace();
            System.exit(2);
            return false;
        }
        return true;
    }

    public void saveScenes(Simulation sim, String baseDir, String simName) {

        // String baseDir = sim.getSessionDir(); //get the name of the simulation's
        // directory
        String sep = System.getProperty("file.separator"); // get the right separator for your operative system
        String currentDir = baseDir + sep;
        BufferedWriter bwout;

        // try {
        //     File currentSimDir = new File(currentDir);
        //     if (!currentSimDir.exists()) {
        //         currentSimDir.mkdirs();
        //     }
        //     sim.saveState(currentDir + currentSim + "_" + simName + ".sim");
        // } catch (Exception e) {
        //     e.printStackTrace();
        //     System.exit(3);
        // }

        try {

            bwout = new BufferedWriter(
                    new FileWriter(resolvePath(simName + "_Report.csv")));
                    Collection<Report> reportCollection = sim.getReportManager().getObjects();

            for (Report thisReport : reportCollection) {
                bwout.write(thisReport.getPresentationName() + ",");
            }


            bwout.write("\n");

            for (Report thisReport : reportCollection) {

                String fieldLocationName = thisReport.getPresentationName();
                Double fieldValue = thisReport.getReportMonitorValue();
                String fieldUnits = thisReport.getUnits().toString();

                // Printing to chek in output window
                sim.println("Field Location :" + fieldLocationName);
                sim.println(" Field Value :" + fieldValue);
                sim.println(" Field Units :" + fieldUnits);
                sim.println("");

                // Write Output file as "sim file name"+report.csv
                bwout.write(fieldValue + ",");

            }

            bwout.close();

            /*for (Scene scn : sim.getSceneManager().getScenes()) {
                sim.println("Saving Scene: " + scn.getPresentationName());
                scn.printAndWait(resolvePath(currentDir + scn.getPresentationName() + ".jpg"), 1, 1920, 1080);
            }*/

            /*for (StarPlot plt : sim.getPlotManager().getObjects()) {
                sim.println("Saving Plot: " + plt.getPresentationName());
                plt.encode(resolvePath(currentDir + plt.getPresentationName() + ".jpg"), "jpg", 1920, 1080);
            }*/

        } catch (IOException iOException) {
            iOException.printStackTrace();
            System.exit(3);
        }

    }

    public void ReadCSVInputs(String fileName) {
        String line = "";
        try {
            BufferedReader br = new BufferedReader(new FileReader(fileName));
            while ((line = br.readLine()) != null) {
                String[] values = line.split(",");
                biplane1AOA = Double.parseDouble(values[0]);
                biplane2AOA = Double.parseDouble(values[1]);
                biplaneGapSize = Double.parseDouble(values[2]);
                biplanePosition = Double.parseDouble(values[3]);
                fw4thAOA = Double.parseDouble(values[4]);
            }
        } catch (IOException e) {
            e.printStackTrace();
            System.exit(1);
        }
    }

    public static void main() { // can use for testing and validating input

    }
}
