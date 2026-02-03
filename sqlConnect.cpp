#include <iostream>
#include <pqxx/pqxx>
#include <string>

int main(){
  try {
    //1.connection to the db
    std::string dbname,user,password {""};

    std::cout<<"enter the database name"<<'\n';
    std::cin>>dbname;
    std::cout<<"enter the username"<<'\n';
    std::cin>>user;
    std::cout<<"enter the username password"<<'\n';
    std::cin>>password;

    std::string prompt = "dbname="+dbname+" "+"user="+user+" "+"password="+password+" "+"host=localhost";

pqxx::connection cx(prompt);
            std::cout << "Connected to " << cx.dbname() << '\n';

            // Start a transaction.  A connection can only have one transaction
            pqxx::work tx{cx};

            for (auto [name, salary] : tx.query<std::string, int>(
                "SELECT name, salary FROM employee ORDER BY name"))
            {
                std::cout << name << " earns " << salary << ".\n";
            }

            for (auto [name, salary] : tx.stream<std::string_view, int>(
                "SELECT name, salary FROM employee"))
            {
                std::cout << name << " earns " << salary << ".\n";
            }

            std::cout << "Doubling all employees' salaries...\n";
            tx.exec("UPDATE employee SET salary = salary*2").no_rows();

            int my_salary = tx.query_value<int>(
                "SELECT salary FROM employee WHERE name = 'Me'");
            std::cout << "I now earn " << my_salary << ".\n";

            auto [top_name, top_salary] = tx.query1<std::string, int>(
                R"(
                    SELECT name, salary
                    FROM employee
                    WHERE salary = (SELECT max(salary) FROM employee)
                    LIMIT 1
                )");
            std::cout << "Top earner is " << top_name << " with a salary of "
                      << top_salary << ".\n";

            pqxx::result res = tx.exec("SELECT * FROM employee");
            std::cout << "Columns:\n";
            for (pqxx::row_size_type col = 0; col < res.columns(); ++col)
                std::cout << res.column_name(col) << '\n';

            std::cout << "Making changes definite: ";
            tx.commit();
            std::cout << "OK.\n";
        }
        catch (std::exception const &e)
        {
            std::cerr << "ERROR: " << e.what() << '\n';
            return 1;
        }
        return 0;
    }









