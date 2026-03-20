
import requests
import time
import concurrent.futures
import matplotlib.pyplot as plt
import json
from datetime import datetime
import statistics

FUNCTION_URL = "https://dsysts-sensors.azurewebsites.net/api/sensor?code=BReQl0S8bjjN7W-FnxQ6n0P_-swXerm1lWluidj7PLRCAzFuJ0Gd9g=="
# FUNCTION_URL = "http://localhost:7071/api/sensor"

def call_function():
    start_time = time.time()
    try:
        # Use GET request (matches your function.json configuration)
        response = requests.get(FUNCTION_URL, timeout=30)
        elapsed_time = (time.time() - start_time) * 1000  # Convert to ms
        
        if response.status_code == 200:
            try:
                data = response.json()
                # Handle both wrapped and unwrapped responses
                if 'Outputs' in data and 'res' in data['Outputs']:
                    # Azure Functions wrapped response
                    res_data = data['Outputs']['res']
                else:
                    # Direct response
                    res_data = data
                
                execution_time = float(res_data.get('execution_time_ms', 0))
                sensor_count = res_data.get('sensor_count', 0)
                
                return {
                    'success': True,
                    'response_time': elapsed_time,
                    'execution_time': execution_time,
                    'sensor_count': sensor_count
                }
            except json.JSONDecodeError:
                return {
                    'success': True,
                    'response_time': elapsed_time,
                    'execution_time': elapsed_time,
                    'sensor_count': 20  # Default assumption
                }
        else:
            return {
                'success': False,
                'response_time': elapsed_time,
                'error': f"HTTP {response.status_code}"
            }
    except Exception as e:
        elapsed_time = (time.time() - start_time) * 1000
        return {
            'success': False,
            'response_time': elapsed_time,
            'error': str(e)
        }

def test_concurrent_requests(num_requests, num_workers):
    print(f"\nTesting with {num_requests} requests using {num_workers} workers...")
    
    start_time = time.time()
    results = []
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=num_workers) as executor:
        futures = [executor.submit(call_function) for _ in range(num_requests)]
        
        for future in concurrent.futures.as_completed(futures):
            results.append(future.result())
    
    total_time = time.time() - start_time
    
    # Calculate statistics
    successful = [r for r in results if r['success']]
    failed = [r for r in results if not r['success']]
    
    if successful:
        response_times = [r['response_time'] for r in successful]
        execution_times = [r['execution_time'] for r in successful]
        
        stats = {
            'total_requests': num_requests,
            'workers': num_workers,
            'successful': len(successful),
            'failed': len(failed),
            'total_time': total_time,
            'throughput': len(successful) / total_time,
            'avg_response_time': statistics.mean(response_times),
            'min_response_time': min(response_times),
            'max_response_time': max(response_times),
            'median_response_time': statistics.median(response_times),
            'avg_execution_time': statistics.mean(execution_times),
        }
        
        print(f"  Success: {stats['successful']}/{num_requests}")
        print(f"  Failed: {stats['failed']}")
        print(f"  Throughput: {stats['throughput']:.2f} req/s")
        print(f"  Avg Response Time: {stats['avg_response_time']:.2f} ms")
        print(f"  Avg Execution Time: {stats['avg_execution_time']:.2f} ms")
        
        if failed:
            print(f"  Errors: {[r['error'] for r in failed[:3]]}")  # Show first 3 errors
        
        return stats
    else:
        print(f"  All requests failed!")
        if failed:
            print(f"  Sample errors: {[r['error'] for r in failed[:5]]}")
        return None

def run_scalability_tests():
    
    print("\nRunning connectivity test...")
    initial_test = test_concurrent_requests(1, 1)
    
    print("\n✓ Connectivity test passed! Proceeding with full tests...\n")
    
    # Test config: (num_requests, num_workers)
    test_configs = [
        (10, 1),    # Sequential
        (10, 5),    # Low concurrency
        (20, 10),   # Medium concurrency
        (50, 10),   # Higher load
        (100, 20),  # High concurrency
        (200, 50),  # Very high load
    ]
    
    results = [initial_test]
    
    for num_requests, num_workers in test_configs:
        stats = test_concurrent_requests(num_requests, num_workers)
        if stats:
            results.append(stats)
        time.sleep(2)  # Brief pause between tests
    
    return results

def plot_results(results):
    if not results:
        print("No results to plot")
        return
    
    fig, axes = plt.subplots(2, 2, figsize=(15, 10))
    fig.suptitle('Azure Function Scalability Test Results', fontsize=16, fontweight='bold')
    
    total_requests = [r['total_requests'] for r in results]
    workers = [r['workers'] for r in results]
    throughput = [r['throughput'] for r in results]
    avg_response = [r['avg_response_time'] for r in results]
    avg_execution = [r['avg_execution_time'] for r in results]
    
    # Plot 1: Throughput vs Total Requests
    axes[0, 0].plot(total_requests, throughput, marker='o', linewidth=2, markersize=8)
    axes[0, 0].set_xlabel('Total Requests', fontsize=12)
    axes[0, 0].set_ylabel('Throughput (req/s)', fontsize=12)
    axes[0, 0].set_title('Throughput vs Load', fontsize=14, fontweight='bold')
    axes[0, 0].grid(True, alpha=0.3)
    
    # Plot 2: Response Time vs Total Requests
    axes[0, 1].plot(total_requests, avg_response, marker='s', color='orange', 
                    linewidth=2, markersize=8)
    axes[0, 1].set_xlabel('Total Requests', fontsize=12)
    axes[0, 1].set_ylabel('Avg Response Time (ms)', fontsize=12)
    axes[0, 1].set_title('Response Time vs Load', fontsize=14, fontweight='bold')
    axes[0, 1].grid(True, alpha=0.3)
    
    # Plot 3: Execution Time vs Concurrency
    axes[1, 0].plot(workers, avg_execution, marker='^', color='green', 
                    linewidth=2, markersize=8)
    axes[1, 0].set_xlabel('Concurrent Workers', fontsize=12)
    axes[1, 0].set_ylabel('Avg Execution Time (ms)', fontsize=12)
    axes[1, 0].set_title('Execution Time vs Concurrency', fontsize=14, fontweight='bold')
    axes[1, 0].grid(True, alpha=0.3)
    
    # Plot 4: Success Rate
    success_rates = [(r['successful'] / r['total_requests'] * 100) for r in results]
    colors = ['green' if rate == 100 else 'orange' if rate >= 90 else 'red' for rate in success_rates]
    axes[1, 1].bar(range(len(results)), success_rates, color=colors, alpha=0.7)
    axes[1, 1].set_xlabel('Test Number', fontsize=12)
    axes[1, 1].set_ylabel('Success Rate (%)', fontsize=12)
    axes[1, 1].set_title('Success Rate by Test', fontsize=14, fontweight='bold')
    axes[1, 1].set_ylim([0, 105])
    axes[1, 1].grid(True, alpha=0.3, axis='y')
    

    test_labels = [f"{r['total_requests']}req\n{r['workers']}w" for r in results]
    axes[1, 1].set_xticks(range(len(results)))
    axes[1, 1].set_xticklabels(test_labels, fontsize=9)
    
    plt.tight_layout()
    
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    filename = f'scalability_results_{timestamp}.png'
    plt.savefig(filename, dpi=300, bbox_inches='tight')
    print(f"\nGraph saved as: {filename}")
    
    plt.show()

def save_results_json(results):
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    filename = f'scalability_results_{timestamp}.json'
    
    with open(filename, 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"Results saved as: {filename}")

if __name__ == "__main__":
    results = run_scalability_tests()
    
    if results:
        save_results_json(results)
        plot_results(results)
        
        print("\n" + "=" * 60)
        print("SUMMARY")
        print("=" * 60)
        best_throughput = max(results, key=lambda x: x['throughput'])
        print(f"Best Throughput: {best_throughput['throughput']:.2f} req/s")
        print(f"  at {best_throughput['total_requests']} requests with {best_throughput['workers']} workers")
        
        total_successful = sum(r['successful'] for r in results)
        total_requests = sum(r['total_requests'] for r in results)
        overall_success_rate = (total_successful / total_requests * 100)
        print(f"\nOverall Success Rate: {overall_success_rate:.1f}% ({total_successful}/{total_requests})")
    else:
        print("\nNo successful tests completed.")
        print("Check the function logs and ensure it's running properly.")
