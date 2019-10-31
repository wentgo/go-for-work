<?php

class JobTest
{
    protected $redis;

    public function __construct()
    {
        $this->redis = new Redis();
        $this->redis->connect('127.0.0.1');
    }

    public function runJob()
    {
         $args = func_get_args();
 
         $name = array_shift($args);
         if (substr($name, -4) != '.php') {
             $name .= '.php';
         }
 
         $job['Name'] = $name;
         $job['Args'] = $args;
 
         $msg = json_encode($job);
         $this->redis->rPush("job:queue", $msg);
    }

    public function delayJob()
    {
         $args = func_get_args();
 
         $time = array_shift($args);
         $name = array_shift($args);
 
         $now = time();
         if ($time < $now) {
             $time += $now;
         }
 
         if (substr($name, -4) != '.php') {
             $name .= '.php';
         }
 
         $job['Time'] = $time;
         $job['Name'] = $name;
         $job['Args'] = $args;
 
         $msg = json_encode($job);
         $this->redis->zAdd("delay:job", $time, $msg);
    }
}

$test = new JobTest();
$test->runJob("Test.php", "Do-It", "Now");
$test->delayJob(5*60, "Test.php", "Do-It", "5-Minute", "Later");

$now = time();
for ($i=0; $i<10; $i++) {
    // Run the job every 10 seconds
    $test->delayJob($i*10, "Test.php", "Job-$i", substr($now, -4));
}
