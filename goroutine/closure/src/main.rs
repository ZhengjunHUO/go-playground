use std::thread;

fn main() {
    let ss = vec!["huo", "foo", "bar"];
    let mut hs = vec![];

    for s in ss {
        let h = thread::spawn(move || {
            println!("{}", s);
        });
        hs.push(h);
    }

    for h in hs {
        h.join().unwrap();
    }
}
